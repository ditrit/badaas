package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	AuthTypeKey string = "auth.type"
)

// AuthType represent the type of authentication
type AuthType string

// The existing AuthTypes
const (
	AuthTypePlain AuthType = "plain"
	AuthTypeOIDC  AuthType = "oidc"
)

// Hold the configuration values for the authentication.
type AuthenticationConfiguration interface {
	ConfigurationHolder
	GetAuthType() AuthType
}

// Concrete implementation of the AuthenticationConfiguration interface
type authenticationConfigurationImpl struct {
	authType AuthType
}

// Instantiate a new configuration holder for the pagination
func NewAuthenticationConfiguration() AuthenticationConfiguration {
	authenticationConfiguration := new(authenticationConfigurationImpl)
	authenticationConfiguration.Reload()
	return authenticationConfiguration
}

// GetAuthType return the auth style to use
func (authenticationConfiguration *authenticationConfigurationImpl) GetAuthType() AuthType {
	return authenticationConfiguration.authType
}

// Reload oidc configuration
func (authenticationConfiguration *authenticationConfigurationImpl) Reload() {
	authenticationConfiguration.authType = AuthType(viper.GetString(AuthTypeKey))
}

// Log the values provided by the configuration holder
func (authenticationConfiguration *authenticationConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("Authentication configuration",
		zap.String("authType", string(authenticationConfiguration.authType)),
	)
}
