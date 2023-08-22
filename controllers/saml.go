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
	"github.com/ditrit/badaas/services/auth/protocols/samlservice"
	"github.com/ditrit/badaas/services/sessionservice"
	"go.uber.org/zap"
)
// On va noter ici les erreurs qui peuvent survenir
var (
	ErrBindingDidNotMatch     = httperrors.NewUnauthorizedError("Only Binding HTTP-Redirect and HTTP-POST are implemented", "Please Restart and use HTTP-Redirect or HTTP-POST as Binding")
	ErrGenerationXML           = httperrors.NewUnauthorizedError("Error during XML generation ", "please recharge the page")
	ErrMissingCookieState        = httperrors.NewUnauthorizedError("authorization code is empty", "please restart auth flow")
	ErrGenerateSLORequest = httperrors.NewUnauthorizedError("Error during decode of the SLO Request", "Error during decode of the SLO Request")
	ErrFormularAnalysis =httperrors.NewUnauthorizedError("Error during formular analysis", http.StatusInternalServerError)
	ErrBuildLogoutResponseBodyFromDocument = httperrors.NewUnauthorizedError("Error during BuildLogoutResponseBodyPostFromDocument", http.StatusInternalServerError)
	ErrWrongMethod =  httperrors.NewUnauthorizedError("Méthode non autorisée seulement la requête post est autorisée", http.StatusMethodNotAllowed) 
	ErrGenerationDocResponseSLO =  httperrors.NewUnauthorizedError("Erreur pendant la génération du document", http.StatusInternalServerError)
	ErrGenerationRedirectURL =  httperrors.NewUnauthorizedError("Erreur pendant la génération de l'URL de redirection", http.StatusInternalServerError)
	ErrBuildBodyForPost = httperrors.NewUnauthorizedError("Erreur pendant la génération du formulaire pour l'authentification", http.StatusInternalServerError)
	ErrNoKeyFound = httperrors.NewUnauthorizedError("The request doesn't contain the key", http.StatusInternalServerError)
	ErrNoUserFindForAuth = httperrors.NewUnauthorizedError("The database doesn't contain the user", http.StatusInternalServerError)
	ErrWrongStatusCode = httperrors.NewUnauthorizedError("There is no Success in the Status Code", http.StatusInternalServerError)
	ErrDecryptandValidateAsser = httperrors.NewUnauthorizedError("Problem with the validation of the Assertion", http.StatusInternalServerError)
	ErrEmptyMail = httperrors.NewUnauthorizedError("No email found in the response", http.StatusInternalServerError)
	ErrUserHasNotContent = httperrors.NewUnauthorizedError("Don't have user with this email", http.StatusInternalServerError)
	ErrWrongBiding= httperrors.NewUnauthorizedError("Wrong Binding for SLO SP initiated", http.StatusInternalServerError)
)



type SAMLController interface {
	SpToIdp(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	IdpToSp(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	BuildSPMetadata(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	HandleLogoutFromIDP(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	HandleLogoutFromSP(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	GenerateMetadataWithSLO(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError)
	
}

/// ajouter le SMAL session
type sAMLController struct {
	logger         *zap.Logger,
	samlService    samlservice.SAMLService,
	sessionservice sessionservice.SessionService,
	userRepository repository.CRUDRepository[models.User, uint],
	samlsession    samlSession,
}

func NewSAMLController(logger *zap.Logger, samlService samlservice.SAMLService,userRepository repository.CRUDRepository[models.User, uint], sessionservice sessionservice.SessionService) SAMLController {
	return &sAMLController{
		logger:         logger,
		samlService:    samlService,
		userRepository: userRepository,
		sessionservice: sessionservice,
		samlsession:    samlSession,
	}
}
func (samlController *sAMLController) HandleLogoutFromSP(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError){
	userID := sAMLController.sessionservice.GetSessionClaimsFromContext(r.Context()).UserID
	for index,uuid :=range sAMLController.samlSession.GetUUIDList(){
		if uuid == userID{
		usernameID := sAMLController.samlSession.GetNameID()[index]
		usersessionIndex := sAMLController.samlSession.GetSessionIndex()[index]
		}
	
	}
	doc,binding,err := samlController.sessionservice.Generate_SLO_Request(usernameID,usersessionIndex)
	if err != nil{
		return nil, ErrGenerateSLORequest
	}
	if binding == "HTTP-Redirect"{
		http.Redirect(w, r, doc , http.StatusFound)
	}
	if binding == "HTTP-POST"{
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w,doc)
	}
	else{
		return nil,ErrWrongBiding
	
	
	} 
	//Faire le code pour utiliser la Méthode POST puis déconnection du user...
	
	return doc,nil


}
func  (samlController *sAMLController) HandleLogoutFromIDP(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError){
	if r.Method != http.MethodPost {
		return nil,ErrWrongMethod // 
	}

	err := r.ParseForm()
	if err != nil {
		return nil,ErrFormularAnalysis 
		
	}

	// Recherche un clé qui contient le terme Request dans la requete post
	var name string 
	for key := range r.PostForm {
		chaineSourceMinuscules := strings.ToLower(key)
		chaineRechercheMinuscules1 := strings.ToLower("Request")
		if strings.Contains(chaineSourceMinuscules, chaineRechercheMinuscules) {
			fmt.Printf("La chaîne '%s' se trouve dans la chaîne source.\n", chaineRecherche)
			name = chaineSourceMinuscules
		} 
	if name==nil || name ==""{
		return nil,ErrNoKeyFound
	}
	var Index string 
	encodedSLORequest := req.Form.Get(name)
	userNameID := samlservice.DecodeIdPSLORequest(encodedRequest)// revoit la requete decoder par le nameID
	// Find index of usernameID to close the session grâce à son uuid et clean la liste
	for  index,nameID := range samlsession.NameID {
		if nameID.value == userNameID {
			Index = index
			break	
		}
	
	}
	if Index != "" {
		userSessionIndex := samlsession.GetSessionIndex()[Index]
		LougoutResponse,err := samlservice.BuildLogoutResponseDocument("Success", userNameID)
		if err == nil {
			return nil, ErrGenerationDocResponseSLO 	
		}
		doc,err := samlservice.BuildLogoutResponseBodyPostFromDocument("",LougoutResponse)
		if err ==nil{
			return nil, ErrBuildLogoutResponseBodyFromDocument
		}
		// trouver comment le faire executer à l'utilisateur y a t'il besoin de faire un truc ... car quand c'est IDP initiated le user nest pas forcément sur e 
		fmt.Fprintf(w,doc)
		sAMLController.sessionService.LogUserOut(sAMLController.sessionservice.GetSessionClaimsFromContext(r.Context()), w)
		samlController.samlsession.CleanSession()
		return doc,nil
	}
	else{
		fmt.Println("User déjà déconnecté")
	}


}

//Send the authentification request 
func (samlController *sAMLController) SpToIdp(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError) {
	if samlServiceProvider.IdentityProviderSSOBinding == "HTTP-Redirect"{
		auth,err :=samlService.BuildRedirectURL()
		if err != nil{
			return nil, ErrGenerationRedirectURL
		}
		http.Redirect(w, r, auth , http.StatusFound)
	}
	if samlServiceProvider.IdentityProviderSSOBinding =="HTTP-POST"{
		body,err := samlService.BuildBodyForPost("")
		if err != nil{
			return nil, ErrBuildBodyForPost
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w,body)
	}
	else:
		return ErrBindingDidNotMatch 


}
// Metadata of the SP can be downloaded in a XML file

func (samlController *sAMLController) GenerateMetadata(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError) {
	metadata :=samlController.samlSession.BuildSPMetadata()
	xmlMeta, err := xml.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return nil,ErrGenerationXML
	}
	if req.Method == http.MethodGet {
		rw.Header().Set("Content-Disposition", "attachment; filename=metadata")
		rw.Header().Set("Content-Type", "application/xml")
		rw.Header().Set("Content-Length", strconv.Itoa(len(xmlMeta)))	
		rw.Write(xmlMeta)
		} 
	else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
	return nil,nil
}
//Methode a utiliser si 
func (samlController *sAMLController) GenerateMetadataWithSLO(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError) {
	metaslo, err := samlController.samlSession.BuildSPMetadataSLO()
		if err != nil {
			fmt.Println("Erreur lors de la génération des métadonnées:", err)
			return nil, ErrGenerationMetadataSLO
		}
		
		xmlMetaslo, err := xml.MarshalIndent(metaslo, "", "    ")
		// Vérifier que la méthode HTTP est GET
		if req.Method == http.MethodGet {
			rw.Header().Set("Content-Disposition", "attachment; filename=metadata")
			rw.Header().Set("Content-Type", "application/xml")
			rw.Header().Set("Content-Length", strconv.Itoa(len(xmlMetaslo)))
			rw.Write(xmlMetaslo)
		} else {
			return nil, ErrWrongMethodForDownload rw.WriteHeader(http.StatusMethodNotAllowed)
		}
}
//Treament of the IDP response only handle post and lead to the creation of a session

func (samlController *sAMLController) IdpToSp(response http.ResponseWriter, request *http.Request) (any, httperrors.HTTPError) {
	err := req.ParseForm()		
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return nil, ErrWrongMethod
	}
	SAML_response_encoded := req.Form.Get("SAMLResponse")
	response_auth, err =  samlService.DecryptandValidateAssertion(SAML_response_encoded)
	if err != nil {
		return nil, ErrDecryptandValidateAsser
	}
	if !strings.Contains(response_auth.Status.StatusCode, "Success"){
		return nil, ErrWrongStatusCode 
	}
	//recherche de toutes les adresse mails appartenant au user 
	MailList := FindEmailAddressesFromResponse(response_auth)
	if len(MailList) == 0 {
		return   nil, ErrEmptyMail
	
	}
	// trouve un email correspondant à l'un de ceux trouvé dans la réponse 
	

	for _,mail := range MailList{
		users, herr := samlController.userRepository.Find(squirrel.Eq{"email": mail}, nil, nil) 
		if users.HasContent{
		break
		}
	}
	fmt.Println("users, herr", users, herr)
	if herr != nil {
		return nil, ErrNoUserFindForAuth
	}

	if !users.HasContent {
		return nil, ErrUserHasNotContent
	}
	///Ajout des crédantials à la session SAML
	for _,Asser := range response_auth.Assertions{
			samlController.samlsession.AddNameID(Asser.Subject.NameID.value)	        
			samlController.samlsession.AddSessionIndex(Asser.AuthnStatement.SessionIndex)
			}
	
	user := users.Ressources[0]
	samlController.samlsession.AddUUIDList(users.Ressources[1])
	samlController.sessionservice.LogUserIn(user, response)
	
	return nil, nil
	
	


}
/////// Implementation de la go routine pour clean la session SAML
go func() {
		for {
			samlController.samlsession.CleanSession()
			time.Sleep(15 * time.Minute) //
		}
	}()
