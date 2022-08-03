package openid_connect

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

//Logger is a middleware handler that does request logging
type Authenticator struct {
	handler http.Handler
}

const prefixAuthMiddleware = "[AUTH MW]"

// This is the list which is used to store user sessions
var AuthenticatedUsers []*User

// This functions checks if the session_code is present in the AuthenticatedUsers list and if the corresponding id_token is valid
func Authorized(codeToVerify string, providerName string) bool {
	for _, u := range AuthenticatedUsers {
		if u.Code == codeToVerify {
			var p Provider = CreateProvider(providerName)
			authenticated := p.Authenticated(u.Tokens.Id_token)
			if authenticated.Value == "true" {
				return true
			}
		}
	}
	return false
}

// Creates a new session for a user. This function returns the session_code for the user.
func NewSessionCode(email string, tokens Tokens) string {
	sessionCode := uuid.New().String()
	u := &User{
		Code:   sessionCode,
		Email:  email,
		Tokens: tokens,
	}
	AuthenticatedUsers = append(AuthenticatedUsers, u)
	fmt.Println("Len(AuthenticatedUsers) : " + strconv.Itoa(len(AuthenticatedUsers)) + "\n")
	return sessionCode
}

// Removes the user session based on his session_code
func RemoveSessionCode(sessionCode string) {
	var temp []*User
	for _, u := range AuthenticatedUsers {
		if u.Code != sessionCode {
			temp = append(temp, u)
		}
	}
	AuthenticatedUsers = temp
	fmt.Println("Len(AuthenticatedUsers) : " + strconv.Itoa(len(AuthenticatedUsers)) + "\n")
}

// Print the Auth Middleware messages
func logInAuthMiddleware(msg string) {
	fmt.Println(prefixAuthMiddleware + " " + msg + "\n")
}

// This middleware checks if the session_code given as a Authorization Bearer header is authorized
func MiddlewareAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		sessionCode := request.Header.Get("Authorization")[7:]

		providers, ok := request.URL.Query()["provider"]
		if !ok || len(providers[0]) < 1 {
			http.Error(w, "Missing provider parameter", http.StatusInternalServerError)
			return
		}
		providerName := providers[0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if Authorized(sessionCode, providerName) {
			logInAuthMiddleware(
				fmt.Sprintf("(ALLOWED) sessionCode %s have been used successfully to get to path %s", sessionCode, request.URL.Path))
			next.ServeHTTP(w, request)
		} else {
			logInAuthMiddleware(
				fmt.Sprintf("(NOT ALLOWED) sessionCode %s have been used unsuccessfully to get to path %s", sessionCode, request.URL.Path))
			w.Write([]byte(
				"{\"status\": \"not authenticated\"}",
			))
		}

	})
}
