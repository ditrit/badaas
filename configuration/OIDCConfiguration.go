package configuration

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// scopeSeparator is the separator used to get a list of string from a single string
	scopeSeparator = ","
)

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
	oidcConfiguration.clientID = viper.GetString(OIDCClientIDKey)
	oidcConfiguration.clientSecret = viper.GetString(OIDCClientSecretIDKey)
	oidcConfiguration.issuer = viper.GetString(OIDCIssuerKey)
	oidcConfiguration.claimIdentifier = viper.GetString(OIDCClaimIdentifierKey)
	oidcConfiguration.redirectURL = viper.GetString(OIDCRedirectURLKey)

	// Get scopes
	scopesStr := viper.GetString(OIDCScopesKey)
	dirtyScopes := strings.Split(scopesStr, scopeSeparator)
	oidcConfiguration.scopes = make([]string, 0)
	for _, dirtyScope := range dirtyScopes {
		cleanScope := strings.TrimSpace(dirtyScope)
		oidcConfiguration.scopes = append(oidcConfiguration.scopes, cleanScope)
	}

}

// Log the values provided by the configuration holder
func (oidcConfiguration *oidcConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("OIDC configuration",
		zap.String("clientID", oidcConfiguration.clientID),
		zap.String("clientSecret", maskPassword(oidcConfiguration.clientSecret)),
		zap.String("issuer", oidcConfiguration.issuer),
		zap.String("redirectUrl", oidcConfiguration.redirectURL),
		zap.String("claimIdentifier", oidcConfiguration.claimIdentifier),
		zap.Strings("scopes", oidcConfiguration.scopes),
	)
}
