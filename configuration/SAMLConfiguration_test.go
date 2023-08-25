package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

var SAMLConfigurationString = `
auth:
 saml:
  TimeValiditySPMetadata: 0
  IdpMetadataFullPath: idpmetadatafullpath
  SAMLSPSSOBindingKey : samlspssobindingkey
  DomaineName: domainename
  SignAuthnRequests: true
  ForceAuthn: false
  SPkeyPath: spkeypath
  SPCertifPath: spcertifpath
  SPSigningKeyPath: spsigningkeypath
  SPSigningCertPath: spsigningcertpath 
`

func TestSAMLConfigurationNewSAMLConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewSAMLConfiguration(), "the contructor for SAMLConfiguration should not return a nil value")

}

func TestSAMLConfigurationGetTimeValiditySPMetadata(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, 0, samlConfiguration.GetTimeValiditySPMetadata())
}

func TestSAMLConfigurationGetIdpMetadataFullPath(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "idpmetadatafullpath", samlConfiguration.GetIdpMetadataFullPath())
}

func TestSAMLConfigurationGetSPKeyPath(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "spkeypath", samlConfiguration.GetSPKeyPath())
}

func TestSAMLConfigurationGetSPCertPath(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "spcertifpath", samlConfiguration.GetSPCertPath())
}

func TestSAMLConfigurationGetSPSigningKeyPath(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "spsigningkeypath", samlConfiguration.GetSPSigningKeyPath())
}

func TestSAMLConfigurationGetSPSigningCertPath(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "spsigningcertpath", samlConfiguration.GetSPSigningCertPath())
}

func TestSAMLConfigurationGetSPSSOBinding(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "samlspssobindingkey", samlConfiguration.GetSPSSOBinding())
}

func TestSAMLConfigurationGetDomaineName(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t,"domainename", samlConfiguration.GetDomaineName())
}

func TestSAMLConfigurationGetSignAuthnRequests(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "true", samlConfiguration.GetSignAuthnRequests())
}

func TestSAMLConfigurationGetForceAuthn(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	samlConfiguration := configuration.NewSAMLConfiguration()
	assert.Equal(t, "false", samlConfiguration.GetForceAuthn())
}
func TestSAMLConfigurationLog(t *testing.T) {
	setupViperEnvironment(SAMLConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	samlConfiguration := configuration.NewSAMLConfiguration()
	samlConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "SAML configuration", log.Message)
	require.Len(t, log.Context, 6)
	// assert.ElementsMatch(t, []zap.Field{
	// 	{Key: "TimeValiditySPMetadata", Type: zapcore.Int64Type, Int64:0 },
	// 	{Key: "IdpMetadataFullPath", Type: zapcore.StringType, String: "idpmetadatafullpath"},
	// 	{Key: "SAMLSPSSOBindingKey", Type: zapcore.StringType, String: "samlspssobindingkey"},
	// 	{Key: "DomaineName", Type: zapcore.StringType, String: "domainename"},
	// 	{Key: "SignAuthnRequests", Type: zapcore.StringType, String: "true"},
	// 	{Key: "ForceAuthn", Type: zapcore.StringType, String: "false"},
	// 	{Key: "SPkeyPath", Type: zapcore.StringType, String: "spkeypath"},
	// 	{Key: "SPCertifPath", Type: zapcore.StringType, String: "spcertifpath"},
	// 	{Key: "SPSigningKeyPath", Type: zapcore.StringType, String: "spsigningkeypath"},
	// 	{Key: "SPSigningCertPath", Type: zapcore.StringType, String: "spsigningcertpath"},
	// }, log.Context)
}
