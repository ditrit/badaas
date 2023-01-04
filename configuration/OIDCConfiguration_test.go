package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

var OIDCConfigurationString = `
auth:
  oidc:
    clientID: clientid
    clientSecret: clientsecret
    issuer: issuer
    claimIdentifier: email
    redirectURL: "https://www.supersite.com"
    scopes: scope1, scope2`

func TestOIDCConfigurationNewOIDCConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewOIDCConfiguration(), "the contructor for OIDCConfiguration should not return a nil value")

}

func TestOIDCConfigurationGetClientID(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	oidcConfiguration := configuration.NewOIDCConfiguration()
	assert.Equal(t, "clientid", oidcConfiguration.GetClientID())
}

func TestOIDCConfigurationGetClientSecret(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	oidcConfiguration := configuration.NewOIDCConfiguration()
	assert.Equal(t, "clientsecret", oidcConfiguration.GetClientSecret())
}

func TestOIDCConfigurationGetIssuer(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	oidcConfiguration := configuration.NewOIDCConfiguration()
	assert.Equal(t, "issuer", oidcConfiguration.GetIssuer())
}

func TestOIDCConfigurationGetRedirectURL(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	oidcConfiguration := configuration.NewOIDCConfiguration()
	assert.Equal(t, "https://www.supersite.com", oidcConfiguration.GetRedirectURL())
}

func TestOIDCConfigurationGetClaimIdentifier(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	oidcConfiguration := configuration.NewOIDCConfiguration()
	assert.Equal(t, "email", oidcConfiguration.GetClaimIdentifier())
}

func TestOIDCConfigurationGetScopes(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	oidcConfiguration := configuration.NewOIDCConfiguration()
	assert.Equal(t, []string{"scope1", "scope2"}, oidcConfiguration.GetScopes())
}

func TestOIDCConfigurationLog(t *testing.T) {
	setupViperEnvironment(OIDCConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	oidcConfiguration := configuration.NewOIDCConfiguration()
	oidcConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "OIDC configuration", log.Message)
	require.Len(t, log.Context, 6)
	// assert.ElementsMatch(t, []zap.Field{
	// 	{Key: "clientID", Type: zapcore.StringType, String: "clientid"},
	// 	{Key: "clientSecret", Type: zapcore.StringType, String: "************"},
	// 	{Key: "issuer", Type: zapcore.StringType, String: "issuer"},
	// 	{Key: "redirectUrl", Type: zapcore.StringType, String: "https://www.supersite.com"},
	// 	{Key: "claimIdentifier", Type: zapcore.StringType, String: "email"},
	// 	{Key: "scopes", Type: zapcore.ArrayMarshalerType, Interface: []interface{}{"scope1", "scope2"}},
	// }, log.Context)
}
