package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/openid_connect"
)

// This controller handles the calls related to the OpenIDConnect flow

// This function redirects the browser of the user to the authentication URL of the requested provider
func LoginScreen(w http.ResponseWriter, request *http.Request) {

	providers, ok := request.URL.Query()["provider"]
	if !ok || len(providers[0]) < 1 {
		http.Error(w, "Missing provider parameter", http.StatusInternalServerError)
		return
	}
	providerName := providers[0]

	log.Println("provider: " + providerName + "\n")

	states, ok := request.URL.Query()["state"]
	if !ok || len(states[0]) < 1 {
		http.Error(w, "Missing state parameter", http.StatusInternalServerError)
		return
	}
	state := states[0]

	nonces, ok := request.URL.Query()["nonce"]
	if !ok || len(nonces[0]) < 1 {
		http.Error(w, "Missing nonce parameter", http.StatusInternalServerError)
		return
	}
	nonce := nonces[0]

	var p openid_connect.OIDCProvider = openid_connect.GetProvider(providerName)

	URL := p.CreateAuthURL(state, nonce)

	http.Redirect(w, request, URL, http.StatusFound)
}

// This function exchanges the OIDC code to get the OIDC tokens given by the provider, then a new authenticated user is created in the backend storage. The session_code corresponding to the new user is sent back to the frontend
func GetSessionCode(w http.ResponseWriter, request *http.Request) {

	providers, ok := request.URL.Query()["provider"]
	if !ok || len(providers[0]) < 1 {
		http.Error(w, "Missing provider parameter", http.StatusInternalServerError)
		return
	}
	providerName := providers[0]

	log.Println("provider: " + providerName + "\n")

	var code openid_connect.Code
	err := json.NewDecoder(request.Body).Decode(&code)
	if err != nil {
		http.Error(w, "Missing code in body json", http.StatusBadRequest)
		return
	}

	log.Println("code: " + code.Value + "\n")

	var p openid_connect.OIDCProvider = openid_connect.GetProvider(providerName)

	tokens, email, nonce, error := p.GetTokens(code.Value)

	sessionCode := openid_connect.NewSessionCode(email, tokens)
	if error != "" {
		http.Error(w, error, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(
			fmt.Sprintf(`{"session_code": "%s","nonce": "%s" }`, sessionCode, nonce),
		))
		return
	}
}

// This function uses the refresh_token to call the RefreshTokens method of the Provider struct. The session of the user is refreshed and the new session_code is sent to the frontend
func RefreshTokens(w http.ResponseWriter, request *http.Request) {

	providers, ok := request.URL.Query()["provider"]
	if !ok || len(providers[0]) < 1 {
		http.Error(w, "Missing provider parameter", http.StatusInternalServerError)
		return
	}
	providerName := providers[0]

	log.Println("provider: " + providerName + "\n")

	sessionCode := request.Header.Get("Authorization")[7:]
	refreshToken := ""
	email := ""

	for _, u := range repository.GetUsers() {
		if u.Code == sessionCode {
			refreshToken = u.Tokens.Refresh_token
			email = u.Email
			break
		}
	}

	log.Println("refreshToken: " + refreshToken + "\n")

	var p openid_connect.OIDCProvider = openid_connect.GetProvider(providerName)

	tokens, error := p.RefreshTokens(refreshToken)

	openid_connect.RemoveSessionCode(sessionCode)

	if error != "" {
		http.Error(w, error, http.StatusInternalServerError)
		return
	} else {
		sessionCode = openid_connect.NewSessionCode(email, tokens)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(
			fmt.Sprintf(`{"session_code": "%s"}`, sessionCode),
		))
	}
}

// This function only sends a json {"status":"authenticated"} as it is only reachable if the session of the user is valid. The checking of the session_code is made in the MiddlewareAuthenticator
func Authenticated(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.Write([]byte(
		"{\"status\": \"authenticated\"}",
	))
}

// This function revokes the refresh_token of the user and deletes the user session from the backend storage
func Logout(w http.ResponseWriter, request *http.Request) {

	providers, ok := request.URL.Query()["provider"]
	if !ok || len(providers[0]) < 1 {
		http.Error(w, "Missing provider parameter", http.StatusInternalServerError)
		return
	}
	providerName := providers[0]

	log.Println("provider: " + providerName + "\n")

	sessionCode := request.Header.Get("Authorization")[7:]

	var p openid_connect.OIDCProvider = openid_connect.GetProvider(providerName)

	for _, u := range repository.GetUsers() {
		if u.Code == sessionCode {
			error := p.RevokeToken(u.Tokens.Refresh_token)
			if error != "" {
				http.Error(w, error, http.StatusInternalServerError)
				return
			}
			openid_connect.RemoveSessionCode(sessionCode)
			break
		}
	}

	w.Write([]byte("Revocation successful"))

}
