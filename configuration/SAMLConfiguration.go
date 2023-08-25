// The config keys regarding the SAML protocol
package configuration

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
	dsigtypes "github.com/russellhaering/goxmldsig/types"
)

const(
	SAMLTimeValiditySPMetadataKey = string "auth.saml.TimeValiditySPMetadata"	
	SAMLIdpMetadataFullPathKey = string "auth.saml.IdpMetadataFullPath" 
	SAMLSPKeyPathKey = string "auth.saml.SPkeyPath"
	SAMLSPCertPathKey = string "auth.saml.SPCertifPath"
	SAMLSPSigningKeyPathKey = string "auth.saml.SPSigningPath"
	SAMLSPSigningCertPathKey = string "auth.saml.SPSigningPath"
	SAMLSPSSOBindingKey = string "auth.saml.SPSSOBinding" 
	SAMLDomaineNameKey = string "auth.saml.DomaineName"
	SAMLSignAuthnRequestsKey = string "auth.saml.SignAuthnRequests"
	SAMLForceAuthnKey = string "auth.saml.ForceAuthn"



	
)

// Hold the configuration values for the saml relying party.
type SAMLConfiguration interface {
	// Be aware that the SAMLConfiguration interface implement the ConfigurationHolder interface.
	ConfigurationHolder 
	GetTimeValiditySPMetadata() int
	
	GetIdpMetadataFullPath() string
	GetSPKeyPath() string
        GetSPCertPath() string
        GetSPSigningKeyPath() string
        GetSPSigningCertPath() string       
        GetSPSSOBinding() string
        GetDomaineName() string 
        GetSignAuthnRequests() bool
        GetForceAuthn() bool

}
var _ SAMLConfiguration = (*samlConfigurationImpl)(nil)
type samlConfigurationImpl struct {
	timeValiditySPMetadata int,
	idpMetadataFullPath,spKeyPath,spCertPath string,
	spSigningKeyPath,spSigningCertPath string,
	spSSOBinding,domaineName string,
	signAuthnRequests,forceAuthn bool,
	
	
}

func NewSAMLConfiguration() SAMLConfiguration {
	samlConfiguration := new(samlConfigurationImpl)
	samlConfiguration.Reload()
	return samlConfiguration
}


func (samlConfiguration *samlConfigurationImpl) GetTimeValiditySPMetadata() int{
	return samlConfiguration.timeValiditySPMetadata
}


func (samlConfiguration *samlConfigurationImpl) GetIdpMetadataFullPath() string{
	return samlConfiguration.idpMetadataFullPath
}

func (samlConfiguration *samlConfigurationImpl) GetSPKeyPath() string{
	return samlConfiguration.spKeyPath
}

func (samlConfiguration *samlConfigurationImpl) GetSPCertPath() string{
	return samlConfiguration.spCertPath
}

func (samlConfiguration *samlConfigurationImpl) GetSPSigningKeyPath() string{
	return samlConfiguration.spSigningKeyPath
}

func (samlConfiguration *samlConfigurationImpl) GetSPSigningCertPath() string{
	return samlConfiguration.spSigningCertPath
}

func (samlConfiguration *samlConfigurationImpl) GetSPSSOBinding() string{
	return samlConfiguration.spSSOBinding
}
/*
func (samlConfiguration *samlConfigurationImpl) GetSPSLOURL() string{
	return samlConfiguration.spSLOURL
}

func (samlConfiguration *samlConfigurationImpl) GetSPIssuer() string{
	return samlConfiguration.spIssuer
}

func (samlConfiguration *samlConfigurationImpl) GetAssertionConsumerServiceURL() string{
	return samlConfiguration.assertionConsumerServiceURL
}

func (samlConfiguration *samlConfigurationImpl) GetAudienceURI() string{
	return samlConfiguration.audienceURI
}
*/

func (samlConfiguration *samlConfigurationImpl) GetDomaineName() string{
	return samlConfiguration.domaineName
}

func (samlConfiguration *samlConfigurationImpl) GetSignAuthnRequests() bool{
	return samlConfiguration.signAuthnRequests
}

func (samlConfiguration *samlConfigurationImpl) GetForceAuthn() bool{
	return samlConfiguration.forceAuthn
}


// Reload the saml configuration
func (samlConfiguration *samlConfigurationImpl) Reload() {

	samlConfiguration.timeValiditySPMetadata = viper.GetInt(SAMLTimeValiditySPMetadataKey)

	samlConfiguration.idpMetadataFullPath = viper.GetString(SAMLIdpMetadataFullPathKey)
	samlConfiguration.spKeyPath = viper.GetString(SAMLSPKeyPathKey)
	samlConfiguration.spCertPath = viper.GetString(SAMLSPCertPathKey)
	samlConfiguration.spSigningKeyPath = viper.GetString(SAMLSPSigningKeyPathKey)
	samlConfiguration.spSigningCertPath = viper.GetString(SAMLSPSigningCertPathKey)
	samlConfiguration.spSSOBinding = viper.GetString(SAMLSPSSOBindingKey)
	samlConfiguration.domaineName = viper.GetBool(SAMLDomaineNameKey)
	samlConfiguration.signAuthnRequests = viper.GetBool(SAMLSignAuthnRequestsKey)
	samlConfiguration.forceAuthn = viper.GetBool(SAMLForceAuthnKey)

	}
	
	
func (samlConfiguration *samlConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("SAML configuration",
		zap.String("TimeValiditySPMetadata", string(samlConfiguration.timeValiditySPMetadata)),
		zap.String("idpMetadataFullPath", samlConfiguration.idpMetadataFullPath),
		zap.String("spKeyPath", samlConfiguration.spKeyPath),
		zap.String("spCertPath", samlConfiguration.spCertPath),
		zap.String("spSigningKeyPath", samlConfiguration.spSigningKeyPath),
		zap.String("spSigningCertPath", samlConfiguration.spSigningCertPath),
		zap.String("spSSOBinding", samlConfiguration.spSSOBinding),
		zap.String("domaineName", samlConfiguration.domaineName),
		zap.String("signAuthnRequests", string(samlConfiguration.signAuthnRequests)),
		zap.String("forceAuthn", string(samlConfiguration.forceAuthn)),		

}


