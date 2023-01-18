# How to handle a new identification protocol in BADAAS. 

## Preamble

Please read the [uber-go/guide](https://github.com/uber-go/guide) to learn how 
golang code should be produced.
We tried to stick to the recommendations of this guide as much as possible.

We will use OIC as an example throughout this guide, but it should be the same 
process with any other indentification protocols.

## 1. Configuration

Create a configuration interface named `OIDCConfiguration` and the implementation.

```go
// The config keys regarding the OIDC protocol
const (
	OIDCClientIDKey        string = "auth.oidc.clientID"
	OIDCClientSecretIDKey  string = "auth.oidc.clientSecret"
	OIDCIssuerKey          string = "auth.oidc.issuer"
	OIDCClaimIdentifierKey string = "auth.oidc.claimIdentifier"
	OIDCRedirectURLKey     string = "auth.oidc.redirectURL"
	OIDCScopesKey          string = "auth.oidc.scopes"
)

// Hold the configuration values for the oidc relying party.
type OIDCConfiguration interface {
	// Be aware that the OIDCConfiguration interface implement the ConfigurationHolder interface.
	ConfigurationHolder 
	GetClientID() string
	GetClientSecret() string
	GetIssuer() string
	GetRedirectURL() string
	GetClaimIdentifier() string
	GetScopes() []string
}

// Concrete implementation of the PaginationConfiguration interface.
type oidcConfigurationImpl struct {
	clientID, clientSecret, issuer, redirectURL, claimIdentifier string
	scopes                                                       []string
}

// Instantiate a new configuration holder for the pagination.
func NewOIDCConfiguration() OIDCConfiguration {
	oidcConfiguration := new(oidcConfigurationImpl)
	oidcConfiguration.Reload()
	return oidcConfiguration
}

// GetRedirectURL return the redirection URL
func (oidcConfiguration *oidcConfigurationImpl) GetRedirectURL() string {
	return oidcConfiguration.redirectURL
}

// GetClientID return the Oauth2 clientID
func (oidcConfiguration *oidcConfigurationImpl) GetClientID() string {
	return oidcConfiguration.clientID
}

// GetClientSecret return the Oauth2 clientSecret
func (oidcConfiguration *oidcConfigurationImpl) GetClientSecret() string {
	return oidcConfiguration.clientSecret
}

// GetIssuer return the Oauth2 issuer
func (oidcConfiguration *oidcConfigurationImpl) GetIssuer() string {
	return oidcConfiguration.issuer
}

// GetClaimIdentifier return the claim that identify the OIDC user
func (oidcConfiguration *oidcConfigurationImpl) GetClaimIdentifier() string {
	return oidcConfiguration.claimIdentifier
}

// GetScopes return the scopes to ask to get the right claims
func (oidcConfiguration *oidcConfigurationImpl) GetScopes() []string {
	return oidcConfiguration.scopes
}

// Reload the oidc configuration
func (oidcConfiguration *oidcConfigurationImpl) Reload() {
	// Reload the configuration values
}

// Log the values provided by the configuration holder
func (oidcConfiguration *oidcConfigurationImpl) Log(logger *zap.Logger) {
	// Log the stored values
}

```

Then we will need to register the config key for the verdeter command.

```go
func initOIDCCommands(cfg *verdeter.VerdeterCommand) {
	cfg.LKey(configuration.OIDCClientIDKey, verdeter.IsStr, "", "The OIDC client ID provided by the OIDC Provider")

	cfg.LKey(configuration.OIDCClientSecretIDKey, verdeter.IsStr, "", "The OIDC client secret provided by the OIDC Provider")

	cfg.LKey(configuration.OIDCIssuerKey, verdeter.IsStr, "", "The OIDC issuer URL (example: https://accounts.google.com)")

	cfg.LKey(configuration.OIDCClaimIdentifierKey, verdeter.IsStr, "",
		"The name of the unique user identifier in the claims of the ID Token returned by the OIDC Provider.")
	cfg.SetDefault(configuration.OIDCClaimIdentifierKey, "sub")

	cfg.LKey(configuration.OIDCRedirectURLKey, verdeter.IsStr, "", "The URL of the callback on the SPA")

	cfg.LKey(configuration.OIDCScopesKey, verdeter.IsStr, "", "The scopes to request to the OIDC Provider.")
}
```

Then we call this function in the in the `init()` function in "commands/rootCmd.go".

## 2. Creating a Service

First let's create an interface named `OIDCService`:

```go 
// OIDCService integrate an OIDC relying party for badaas.
type OIDCService interface {
	BuildRedirectURL(state string) string
	ExchangeAuthorizationCode(code string) (userIdentifier string, err error)
}
```

This interface is intentionally minimal, it is only used by the HTTP controller
we don't want to have to support many methods (especially if we don't use them).

Then we will create a struct :

```go
// oidcService is a concrete implementation of OIDCService.
type oidcService struct {
	logger *zap.Logger
	// oauth2Config oauth2
	oauth2Config oauth2.Config
	// The provider OIDC
	provider *oidc.Provider

	// configuration
	oidConfiguration configuration.OIDCConfiguration
}
```

Obviously the methods have to be implemented.

## 3. Create a http controller

Then we create a HTTP controller in the `controllers` package.

```go

var (
	ErrStateDidNotMatch     = httperrors.NewUnauthorizedError("state did not match", "please restart auth flow")
	ErrFailedToEchangeToken = func(err error) httperrors.HTTPError {
		return httperrors.NewInternalServerError("failed to exchange token", "please restart auth flow", err)
	}
	ErrAuthorizationCodeNotFound = httperrors.NewUnauthorizedError("authorization code is empty", "please restart auth flow")
	ErrMissingCookieState        = httperrors.NewHTTPError(http.StatusBadRequest, "missing cookie state", "", nil, true)
)

// OIDCController handle http requests for badaas OIDC relying party.
type OIDCController interface {
	RedirectURL(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	CallBack(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
}

// This controller handles the calls related to the OIDC flow.
type oIDCController struct {
	logger         *zap.Logger
	oidcService    oidcservice.OIDCService
	sessionservice sessionservice.SessionService
	userRepository repository.CRUDRepository[models.User, uint]
}

// NewOIDCController is the constructor for OIDCController.
func NewOIDCController(logger *zap.Logger, oidcService oidcservice.OIDCService,
	userRepository repository.CRUDRepository[models.User, uint], sessionservice sessionservice.SessionService) OIDCController {
	return &oIDCController{
		logger:         logger,
		oidcService:    oidcService,
		userRepository: userRepository,
		sessionservice: sessionservice,
	}
}

// RedirectURL return a OIDC redirect URL.
func (oidcController *oIDCController) RedirectURL(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError) {
	state := randString(16)
	setCallbackCookie(response, request, "state", state)

	URL := oidcController.oidcService.BuildRedirectURL(state)

	return dto.OIDCRedirectURL{RedirectURL: URL}, nil
}

// CallBack handle the return of the user from the OIDC provider authentication portal.
func (oidcController *oIDCController) CallBack(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError) {
	state, err := request.Cookie("state")
	if err != nil {
		return nil, ErrMissingCookieState
	}

	if request.URL.Query().Get("state") != state.Value {
		return nil, ErrStateDidNotMatch
	}

	authorizationCode := request.URL.Query().Get("code")
	if authorizationCode == "" {
		return nil, ErrAuthorizationCodeNotFound
	}

	oidcIdentifier, err := oidcController.oidcService.ExchangeAuthorizationCode(authorizationCode)
	if err != nil {
		return nil, ErrFailedToEchangeToken(err)
	}

	users, herr := oidcController.userRepository.Find(squirrel.Eq{"oidc_identifier": oidcIdentifier}, nil, nil)
	fmt.Println("users, herr", users, herr)
	if herr != nil {
		return nil, herr
	}

	if !users.HasContent {
		return nil, httperrors.NewErrorNotFound("user",
			fmt.Sprintf("no user found with oidcIdentifier %q", oidcIdentifier))
	}

	user := users.Ressources[0]

	oidcController.sessionservice.LogUserIn(user, response)
	return nil, nil
}

// Set the callback cookie.
func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	...
}

// create the callback cookie.
func createCallbackCookie(name string, value string, r *http.Request) *http.Cookie {
	...
}

// randString create a random string of nChar.
func randString(nChar int) string {
	...
}
```

Please note that we return 2 interesting types: a payload that will be passed 
to json.Marshall and an HTTPError(the httperrors package implement some contructors).

## 4. Then add the routes

Add the routes to the router package. Take exemple on the other routes.
