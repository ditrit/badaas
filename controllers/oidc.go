package controllers

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/auth/protocols/oidcservice"
	"github.com/ditrit/badaas/services/sessionservice"
	"go.uber.org/zap"
)

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
	c := createCallbackCookie(name, value, r)
	http.SetCookie(w, c)
}

// create the callback cookie.
func createCallbackCookie(name string, value string, r *http.Request) *http.Cookie {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int((time.Minute * 15).Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
		Path:     "/",
	}
	return c
}

// randString create a random string of nChar.
func randString(nChar int) string {
	const (
		allowedCharsInState = "azertyuiopmlkjhgfdsqwxcvbnAZERTYUIOPMLKJHGFDSQWXCVBN0123456789"
	)
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, nChar)
	for i := 0; i < nChar; i++ {
		result[i] = allowedCharsInState[rand.Intn(len(allowedCharsInState))]
	}
	return base64.RawURLEncoding.EncodeToString(result)[:nChar] // strip the one extra byte we get from half the results.
}
