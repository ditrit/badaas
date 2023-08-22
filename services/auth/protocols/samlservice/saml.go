package samlservice
//Le service SAML actuelle ne prend en charge que les bindings POST du idp et du SP
//Il serait bon de configurer le SP pour qu'il puisse envoyer un formulaire POST à l'IDP
// Faire un code qui gère le logout

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/services/sessionservice"
        "crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/base64"
	"encoding/xml"
	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig"github.com/russellhaering/goxmldsig"
)
const(
	SPIssuer = "/auth/saml/authentication"
	SPAssertionConsumtion ="/auth/saml/acs"
	SPMetadataSLO = "/auth/saml/metadataslo"
	SPSLOAddressForIdp = "/auth/saml/sloFromOutSide"

)
var (
	// ErrFailedToExchangeAuthorizationCode is returned when the oauth2Config
	// fail to exchange the tokens against the authorization code.
	ErrFailedToBuildRedictURL = errors.New("Failed to build the redirectURL")
 	ErrFailedDecodeSLORequest = errors.New("Failed to decode SLO Request")
 	ErrFailedDecodeSLOResponse = errors.New("Failed to decode SLO Response")
 	ErrFailedGenerateSLORequest = errors.New("Failed to generate SLO Request")
 	ErrFailedGenerateSLOResponse = errors.New("Failed to generate SLO Reponse")
 	ErrFailedBuildBodyForPost = errors.New("Failed to build body for POST")
 	ErrFailedToBuildMetadataSLO =errors.New("Failed to generate Metadata with SLO Binding")
 	ErrFailedDecodeIdPSLORequest = errors.New("Failed to Decode IdP SLO Request")
 	ErrFailedDecryptandValidateAssertion = errors.New("Failed to Decrypt and Validate Assertion")
 	ErrFailedValidateEncodedResponse = errors.New("Failed Validate Encoded Response")
 	ErrMailListEmpty = errors.New("MailList is empty no email corresponding to the user was found in the service")
	// ErrClaimIdentifierNotFoundInIDTokenBody is returned when the OIDC claim is not found in the ID token.
	ErrFailedToBuildMetadata = errors.New("Failed to generate Metadata")
	ErrFailedLoadX509KeyPair =  errors.New("Failed to load SP certificates and/or key")
	ErrFailedReadIdpMetadataFile = errors.New("Failed to read IDP Metadata File")
	ErrFailedXMLGeneration = errors.New("Failed to generate the XML from the metadata")
	ErrX509ParseCertificate = errors.New("Failed to parse data X509 Certificate")
	ErrBase64DecodingX509Certif = errors.New("Failed to decode base64 from X509 certificate (IDP)")
	ErrMetadataEmpty = errors.New("IDP's metadatas are empty...")
	ErrWrongBinding = errors.New("You probably choose a wrong spbinding please set SAMLSPSSOBindingKey with HTTP-Redirect or HTTP-POST") 
	ErrIDPDontHandleSLOBinding = errors.New("The IDP's metadata don't have any URL which can handle SLO with HTTP-POST binding")
	ErrGenerateSLORequest = errors.New("Problem during generation of the SLO request")
)

type SAMLService interface {
	BuildRedirectURL(RelayState string) string
	BuildSPMetadata() string 
	BuildSPMetadataSLO() string 
	DecryptandValidateAssertion() string
	GetID (string decryptAssertion) [] string 
	Generate_SLO_Request( nameID string, sessionIndex string, binding string)  (*etree.Document, error)
	Decode_SLO_Request(encodedRequest string) (*types.LogoutRequest, error)
	Decode_SLO_Response(encodedRequest string)  (*types.LogoutResponse, error)
}

var _ SAMLService = (*samlService)(nil)

type samlService struct {
	logger *zap.Logger
	//On set le service provider dit SP 
	samlServiceProvider  *saml2.SAMLServiceProvider 
	samlConfiguration *configuration.SAMLConfiguration
}
// Fonctions for the logout 
func (samlService *samlService) Decode_SLO_Request(encodedRequest string) (*types.LogoutRequest, error){
	LogoutReq,err := samlService.samlServiceProvider.ValidateEncodedLogoutRequestPOST(encodedRequest)
	if err != nil {
		return nil, ErrFailedDecodeSLORequest
	}
	return LogoutReq,err
}
func (samlService *samlService) Decode_SLO_Response(encodedRequest string)  (*types.LogoutResponse, error){
	LogoutReq,err := samlService.samlServiceProvider.ValidateEncodedLogoutRequestPOST(encodedRequest)
	if err != nil {
		return nil,ErrFailedDecodeSLOResponse
	}
	return LogoutReq,err
}

func (samlService *samlService) Generate_SLO_Request( nameID string, sessionIndex string)  (string,string,error) {
	SLO_Request,err := samlService.samlServiceProvider.BuildLogoutRequestDocument(nameID string, sessionIndex string)
	
	if err != nil {
		return nil,ErrFailedGenerateSLORequest
	}
	if samlConfiguration.GetSPSSOBinding() =="HTTP-POST"{
		SLOReq,err :=  samlService.samlServiceProvider.BuildLogoutBodyPostFromDocument("",SLO_Request)
		SLOReqString := string(SLOReq)
	}
	if samlConfiguration.GetSPSSOBinding() =="HTTP-Redirect"{
		SLOReqString,err :=  samlService.samlServiceProviderBuildLogoutURLRedirect("",SLO_Request)
		
	}
	if err != nil{
		return nil, ErrGenerateSLORequest
	}
	return SLOReqString,samlConfiguration.GetSPSSOBinding(),err
}

func (samlService *samlService) Generate_SLO_Response(statusCodeValue string, reqID string)  (*etree.Document, error) {
	SLO_Response,err := samlService.samlServiceProvider.buildLogoutResponse(statusCodeValue string, reqID string, true)
	if err != nil {
		return nil, ErrFailedGenerateSLOResponse
	}
	return SLO_Response,err
}
// End of logout function
func (samlService *samlService) BuildRedirectURL() (string,error){
      authURL, err := samlService.samlServiceProvider.BuildAuthURL("")
      if err != nil{
		return "", ErrFailedToBuildRedictURL
		}
      return authURL, err

}
//Cette fonction sert à générer le corps d'un requête post pour authentification
func (samlService *samlService)	BuildBodyForPost() (string,error){

      Body, err := samlService.samlServiceProvider.BuildAuthBodyPost("")
      if err != nil{
		return "", ErrFailedBuildBodyForPost
		}
      return string(Body), err

}
func (samlService *samlService) BuildSPMetadata() (string,error) {
	metadata,err := samlService.samlServiceProvider.Metadata()//Attention les metadatas générer peuvent être considérer comme invalide pour certains IdP à causse de la valeur du champ "validUntil" qui contient trop de chiffre après la virgule garder juste 3 ou effacer les tous. 
	if err != nil {
		return "", ErrFailedToBuildMetadata
		
	}
	return metadata
}

func (samlService *samlService) BuildSPMetadataSLO() (string,error) {
	metadata,err := samlService.samlServiceProvider.MetadataWithSLO(samlConfiguration.GetTimeValiditySPMetadata())//Attention les metadatas générer peuvent être considérer comme invalide pour certains IdP à causse de la valeur du champ "validUntil" qui contient trop de chiffre après la virgule garder juste 3 ou effacer les tous. 
	if err != nil {
		return "", ErrFailedToBuildMetadataSLO
		
	}
	return metadata,err
}

func (samlController *sAMLController) DecodeIdPSLORequest(encodedRequest) (string,error) {
		decodedSLORequest, err := samlService.ValidateEncodedLogoutRequestPOST(encodedRequest)
		if err != nil {
			return "", ErrFailedDecodeIdPSLORequest
		
		}
		return decodedSLORequest.NameID
	}


func (samlService *samlService) DecryptandValidateAssertion(Assertion string) (string, error) {
	var response_auth *types.Response
	response_auth, err = sp.ValidateEncodedResponse(Assertion)
	if err != nil {
		return "",ErrFailedDecryptandValidateAssertion
	}
	return reponse_auth
}

func (samlService *samlService) FindEmailAddressesFromResponse(response_auth *types.Response) ([]string, error) {
	var MailList []string
		response_auth, err = sp.ValidateEncodedResponse(SAML_response_encoded)
		if err != nil {
			return nil,ErrFailedValidateEncodedResponse
		}
		for _,Asser := range response_auth.Assertions{
			//fmt.Println(Asser.AttributeStatement)
			for _,Att := range Asser.AttributeStatement.Attributes{
				
				for  _,Vals := range Att.Values{
				
					if strings.Contains(Vals.Value, email_detector){  // Rajouter le cas dans la réponse d'avoir le terme success
				
						MailList = append(MailList, Vals.Value)
						
					}
				
				}

				
			}
			
		}
		if len(MailList) == 0 {
			return nil, ErrMailListEmpty
		}
	
	return 	MailList,err	
}


func NewSAMLService(logger *zap.Logger, samlConfiguration configuration.SAMLConfiguration) (SAMLService, error) {
	ctx := context.Background()
//Bizarre le code ...
	provider, err := saml.NewProvider(ctx, samlConfiguration.GetIssuer())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize saml provider with issuer %q, error=%q",
			samlConfiguration.GetIssuer(), err.Error())
	}
	
//On charge les clés du SP 
	key,err := tls.LoadX509KeyPair(samlConfiguration.spKeyPath, samlConfiguration.spCertPath)
	 if err != nil {
		return nil, ErrFailedLoadX509KeyPair
	}
	Key := dsig.TLSCertKeyStore(key)
//On charge les métadatas de l'idp
	rawMetadata, err := ioutil.ReadFile(samlConfiguration.idpMetadataFullPath) // lire le fichier qui nous permet de trust l'idp !! enlever les sauts à la ligne !
        if err != nil {
		return nil, ErrFailedReadIdpMetadataFile
        }
	rawMetadata_clean := []byte(strings.ReplaceAll(string(rawMetadata),"\n",""))
	metadata := &types.EntityDescriptor{
	}
	err = xml.Unmarshal(rawMetadata_clean, metadata)
	if err != nil {
		return nil, ErrFailedXMLGeneration
	}
	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}
	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
			
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
	
				return nil, ErrMetadataEmpty
			}
			fmt.Println()
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				return nil, ErrBase64DecodingX509Certif
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {

				return nil, ErrX509ParseCertificate
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}
	var IDPSSOURL string
	for index,SSOURL := range metadata.IDPSSODescriptor.SingleSignOnService{
		if strings.Contains(SSOURL.Binding,samlConfiguration.spSSOBinding){
			IDPSSOURL = SSOURL

		}
		
	}
	if IDPSSOURL == nil || IDPSSOURL = "" {
		return nil, ErrWrongBinding//Wrong binding
		
	}
	var IDPSLOURL string	
	for index,SLOURL := range metadata.IDPSSODescriptor.SingleLogoutService{
		if strings.Contains(SLOURL.Binding,"HTTP-POST"){
			IDPSLOURL := SLOURL

		}
		
	}
	if IDPSLOURL == nil || IDPSLOURL = "" {
		return nil,ErrIDPDontHandleSLOBinding//Wrong IDP doesn't have a URL which can handle HTTP-POST
		
	}
	&saml2.SAMLServiceProvider{
	
		IdentityProviderSSOURL:      IDPSSOURL,
		IdentityProviderIssuer:      metadata.EntityID,
		//IdentityProviderSSOBinding:  samlConfiguration.GetIdpSSOBinding(), // Only handle redirect binding now have to implement POST binding 
		
		IdentityProviderSLOURL:      IDPSLOURL, 
	      //  IdentityProviderSLOBinding:  samlConfiguration.GetIdpSLOBinding(),
		
		ServiceProviderIssuer:       samlConfiguration.GetDomaineName()+SPIssuer, // ressource qui redirige
		ServiceProviderSLOURL:       samlConfiguration.GetDomaineName()+SPSLOAddressForIdp,
		AssertionConsumerServiceURL: samlConfiguration.GetDomaineName()+SPAssertionConsumtion, //location where the assertion of the Idp will be proceed by the SP
		
		SignAuthnRequests:           samlConfiguration.GetSignAuthnRequests(),
		
		
		ForceAuthn:    		     samlConfiguration.GetForceAuthn(), // If true the idp will authenticate the user even if this user is already log in this idp
		IsPassive:     		     false, 
		
		AudienceURI:                 samlConfiguration.GetDomaineName()+SPIssuer,//It is the SP ID for the IdP usually equal to ServiceProviderIssuer
		
		
		SPKeyStore:                  Key ,//Key public and private
		//SPSigningKeyStore       samlConfiguration.GetSPSigningKey(), // Optional signing key
		IDPCertificateStore:         &certStore, //Public Key of Idp
	}
	return &samlService{logger, config, provider, samlConfiguration}, nil
}
// Ajout lier à la création d'une session SAML 
// Penser à gérer les erreurs


func removeElementAtIndex(slice []interface{}, index int) []int {
	if index < 0 || index >= len(slice) {
		// Index out of bounds, return the original slice unchanged
		return slice
	}

	// Use append to remove the element at the specified index
	// by slicing the slice before and after the index
	return append(slice[:index], slice[index+1:])
}

type SAMLSession interface {// faire un heritage de la session de Badaas

	GetSessionIndex() []string
	GetNameID() []string
		
	AddUUIDList(UserUUID sessionservice.UserID) nil
	AddSessionIndex(UserSessionIndex string) nil 
	AddNameID(UserNameID string) nil
	
	SuppElementFromUUIDList(index) nil 
	SuppElementFromSessionIndex(index) nil	
	SuppElementFromNameID(index) nil
	

	CleanSession() nil
	
}
var _ SAMLSession = (*samlSession)(nil)

type samlSession struct {
	sessionservice sessionservice.SessionService,
	UUIDList []sessionservice.UserID, 
	SessionIndex []string,
	NameID []string,

}
func  NewsamlSession(UUIDList []sessionservice.UserID, SessionIndex []string , NameID []string){
	return &samlSession{
		UUIDList []sessionservice.UserID, 
		SessionIndex []string, 
		NameID []string,
		sessionservice: sessionservice,
	}

}

func (samlSession *samlSession) GetSessionIndex() []string{
	return samlSession.SessionIndex
}

func (samlSession *samlSession) GetNameID() []string{
	return samlSession.NameID
}

func (samlSession *samlSession) GetUUIDList() []sessionservice.UserID  {
	return samlSession.UUIDList
}

func (samlSession *samlSession) AddNameID(UserNameID string)  {
	samlSession.NameID := append(samlSession.GetNameID(),UserNameID)

}

func (samlSession *samlSession) AddSessionIndex(UserSessionIndex string) {
	samlSession.SessionIndex := append(samlSession.GetSessionIndex(),UserSessionIndex)

}

func (samlSession *samlSession) AddUUIDList(UserUUID sessionservice.UserID) {
	samlSession.SessionIndex := append(samlSession.GetUUIDList(), UserUUID)

}
//
func (samlSession *samlSession) SuppElementFromUUIDList(index)  {
	samlSession.UUIDList := removeElementAtIndex(samlSession.GetUUIDList(), index)

}

func (samlSession *samlSession) SuppElementFromSessionIndex(index) {
	samlSession.SessionIndex := removeElementAtIndex(samlSession.GetSessionIndex(), index)

}

func (samlSession *samlSession) SuppElementFromNameID(index) {
	samlSession.NameID := removeElementAtIndex(samlSession.GetNameID(), index)

}



func (samlSession *samlSession) CleanSession() {
	for index,uuid := range samlSession.GetUUIDList(){
	
		if !samlSession.sessionservice.IsValid(uuid)[0] {
			samlSession.SuppElementFromUUIDList(index)
			samlSession.SuppElementFromSessionIndex(index)
			samlSession.SuppElementFromNameID(index)
		}
		
	}

}

// Fin de l'ajout
