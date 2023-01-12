package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var AuthenticationConfigurationString = `
auth:
  type: plain`

func TestAuthenticationConfigurationNewAuthenticationConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewAuthenticationConfiguration(), "the contructor for AuthenticationConfiguration should not return a nil value")

}

func TestAuthenticationConfigurationGetAuthType(t *testing.T) {
	setupViperEnvironment(AuthenticationConfigurationString)
	authentificationConfiguration := configuration.NewAuthenticationConfiguration()
	assert.Equal(t, configuration.AuthTypePlain, authentificationConfiguration.GetAuthType())
}

func TestAuthenticationConfigurationLog(t *testing.T) {
	setupViperEnvironment(AuthenticationConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	authentificationConfiguration := configuration.NewAuthenticationConfiguration()
	authentificationConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Authentication configuration", log.Message)
	require.Len(t, log.Context, 1)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "authType", Type: zapcore.StringType, String: "plain"},
	}, log.Context)
}
