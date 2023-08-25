package main

import (
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

	"os"

	"strings"


)

const (
	email_detector = "@"
)
var NameID string
var SessionIndex string
func main() {
//le code permet de fetch directment les metadata mais probleme de format
	/*res, err := http.Get("https://samltest.id/saml/idp")
	if err != nil {
		panic(err)
	}

	rawMetadata, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(rawMetadata))
		
	if err != nil {
		panic(err)
	}
        println("1")
        
        
       
        ICI on charge les metadata de l'IdP pour instaurer la relations de confiance
        */
        rawMetadata, err := ioutil.ReadFile("idp") // lire le fichier qui nous permet de trust l'idp !! enlever les sauts à la ligne !
        if err != nil {
        fmt.Println(err)
        }
	rawMetadata_clean := []byte(strings.ReplaceAll(string(rawMetadata),"\n",""))
	metadata := &types.EntityDescriptor{
	}
	if metadata == nil {
    		fmt.Println("Le pointeur est nul.")
	} else {
    		fmt.Println("Le pointeur n'est pas nul.")
	}
	err = xml.Unmarshal(rawMetadata_clean, metadata)
	if err != nil {
		panic(err)
	}
	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}
	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
			
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
	
				panic(fmt.Errorf("metadata certificate(%d) must not be empty", idx))
			}
			fmt.Println()
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				fmt.Println(xcert.Data)
				fmt.Println(err)
				panic(err)
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {

				panic(err)
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}
        
        // Fin du code qui traite les metadata de l'IdP
	
	key,err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	 if err != nil {
		panic(err) 
	}
	Key := dsig.TLSCertKeyStore(key)
	
	
	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:   "https://samltest.id/idp/profile/SAML2/Redirect/SSO",// metadata.IDPSSODescriptor.SingleSignOnServices[3].Location,
		IdentityProviderSLOURL:   "https://samltest.id/idp/profile/SAML2/Redirect/SLO",
		ServiceProviderSLOURL:    "http://localhost:80/lougoutBinding",
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       "http://localhost:80/hello", // resource qui redirige / à protéger
		AssertionConsumerServiceURL: "http://localhost:80/v1/_saml_callback",//location to send the POST form (SAML assertion)
		SignAuthnRequests:           false,
		AudienceURI:                 "http://localhost:80/hello",// est un identifiant 
		SPKeyStore:                  Key,
		IDPCertificateStore:         &certStore,
	}

	meta, err := sp.Metadata() 
	if err != nil {
		fmt.Println("Erreur lors de la génération des métadonnées:", err)
		panic(err)
	}
	xmlMeta, err := xml.MarshalIndent(meta, "", "    ")

	// Permet de généner les métadatas rendu téléchargable la variable Until est trop précise enlever les chiffres après la virgule avec un éditeur de texte
	http.HandleFunc("/metadata", func(rw http.ResponseWriter, req *http.Request) {
		// Vérifier que la méthode HTTP est GET
		if req.Method == http.MethodGet {
			rw.Header().Set("Content-Disposition", "attachment; filename=metadata")
			rw.Header().Set("Content-Type", "application/xml")
			rw.Header().Set("Content-Length", strconv.Itoa(len(xmlMeta)))
			rw.Write(xmlMeta)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
// Permet de générer les métadatas rendu 

http.HandleFunc("/metadataslo", func(rw http.ResponseWriter, req *http.Request) {
		var validityHours int64 =0
		metaslo, err := sp.MetadataWithSLO(validityHours) ///Attention il faut 
		if err != nil {
			fmt.Println("Erreur lors de la génération des métadonnées:", err)
			panic(err)
		}
		
		xmlMetaslo, err := xml.MarshalIndent(metaslo, "", "    ")
		// Vérifier que la méthode HTTP est GET
		if req.Method == http.MethodGet {
			rw.Header().Set("Content-Disposition", "attachment; filename=metadata")
			rw.Header().Set("Content-Type", "application/xml")
			rw.Header().Set("Content-Length", strconv.Itoa(len(xmlMetaslo)))
			rw.Write(xmlMetaslo)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	})


http.HandleFunc("/SLO", func(rw http.ResponseWriter, req *http.Request) {
		validityHours = 365*24
		metaslo, err := sp.MetadataWithSLO(validityHours) ///Attention 
		xmlMetaslo, err := xml.MarshalIndent(metaslo, "", "    ")
		// Vérifier que la méthode HTTP est GET
		if req.Method == http.MethodGet {
			rw.Header().Set("Content-Disposition", "attachment; filename=metadata")
			rw.Header().Set("Content-Type", "application/xml")
			rw.Write(xmlMetaslo)
		} else{
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	})


        http.HandleFunc("/v1/_saml_callback", func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()	
		println("1")
		
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// Ce code permet de voir les différents champs dans la requête post que l'idp envoit.
		fmt.Println("Clés de la requête POST:")
		for key := range req.Form {
			fmt.Println(key)
		}

		
		fmt.Printf("Méthode : %s\n", req.Method)
		fmt.Printf("URL : %s\n", req.URL.String())
		fmt.Printf("En-têtes : %v\n", req.Header)
		fmt.Printf("Corps : %v\n", req.Body)
		fmt.Printf("longueur : %v\n", req.ContentLength)
		SAML_response_encoded := req.Form.Get("SAMLResponse")
		fmt.Printf("SAML_response_encoded %v\n",SAML_response_encoded)
		SAML_response_bytes, err := base64.StdEncoding.DecodeString(SAML_response_encoded)
		if err != nil {
			fmt.Println("Erreur lors du décodage en base64:", err)
			return
		}
		SAML_response := string(SAML_response_bytes) //On a le doc XML
		fmt.Printf("%v\n",SAML_response)
	
		
		var response_auth *types.Response
		var MailList []string
		response_auth, err = sp.ValidateEncodedResponse(SAML_response_encoded)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println(response_auth.Status.StatusCode)
		for _,Asser := range response_auth.Assertions{
		        NameID = string(Asser.Subject.NameID.Value)
		        SessionIndex = Asser.AuthnStatement.SessionIndex
			//fmt.Println("ID",Asser.Subject.NameID)
			//fmt.Println("session",Asser.AuthnStatement.SessionIndex)
			
			for _,Att := range Asser.AttributeStatement.Attributes{
				
				for  _,Vals := range Att.Values{
				
					if strings.Contains(Vals.Value, email_detector){
				
						MailList = append(MailList, Vals.Value)
						
					}
				
				}

				
			}
			
		}

	})
	http.HandleFunc("/lougoutBinding", func(w http.ResponseWriter, r *http.Request) {
	
	SLO_Request,err := sp.BuildLogoutRequestDocument(NameID, SessionIndex)
	fmt.Printf("NameID", NameID)
	fmt.Printf("SessionIndex", SessionIndex)
	fmt.Printf("DocumentSLO",SLO_Request)
	SLO_Request.Indent(2)
	SLO_Request.WriteTo(os.Stdout)
/*	if err != nil {
		fmt.Printf("err Build SLO Doc")
	}
*/	
//	SLOReq,err :=  sp.BuildLogoutBodyPostFromDocument("",SLO_Request)
//	SLOReqString := string(SLOReq)


	doc,err :=  sp.BuildLogoutURLRedirect("",SLO_Request)
	if err != nil {
		fmt.Printf("err Build SLO Doc")
	}	



	http.Redirect(w, r, doc , http.StatusFound)
	//w.Header().Set("Content-Type", "text/html")
	//fmt.Fprintf(w,SLOReqString)

	
				
	})
	

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		println("Visit this URL To Authenticate:")
		authURL, err := sp.BuildAuthURL("")
	     	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête POST:", err)
		return
	}	
	http.Redirect(w, r, authURL , http.StatusFound)
		/*
		///HTTP-POST Binding 
		body,err := sp.BuildAuthBodyPost("/")
		if err != nil {
				println("ALERT")
			return
		}
		println("3")
		w.Header().Set("Content-Type", "text/html")
		println(string(body))
		fmt.Fprintf(w,string(body))
		*/

	/*

	        formData := url.Values{
		"SAMLRequest": {base64.StdEncoding.EncodeToString(body)}, // Encodage base64 de la demande SAML
	}
			resp, err := http.PostForm("https://samltest.id/idp/profile/SAML2/POST/SSO", formData)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête POST:", err)
		return
	}
	defer resp.Body.Close()
	//Redirect binding
	        authURL, err := sp.BuildAuthURL("")
	     	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête POST:", err)
		return
	}	
	//	http.Redirect(w, r, authURL , http.StatusFound)
	*/
	})
	
	
	println("Supply:")
	fmt.Printf("test",sp.IdentityProviderSSOURL)
	fmt.Printf("SP ACS URL      : %s\n", sp.AssertionConsumerServiceURL)
		
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
	


	
}

