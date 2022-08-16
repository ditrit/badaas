package openid_connect

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/google/uuid"
)

const prefixAuthMiddleware = "[AUTH MW]"

// This functions checks if the session_code is present in the list of users and if the corresponding id_token is valid
func Authorized(codeToVerify string, providerName string) bool {
	for _, u := range repository.GetUsers() {
		if u.Code == codeToVerify {
			var p OIDCProvider = CreateProvider(providerName)
			authenticated := p.Authenticated(u.Tokens.Id_token)
			if authenticated.Value == "true" {
				return true
			}
		}
	}
	return false
}

// Creates a new session for a user. This function returns the session_code for the user.
func NewSessionCode(email string, tokens models.Tokens) string {
	sessionCode := uuid.New().String()
	u := &models.User{
		Code:   sessionCode,
		Email:  email,
		Tokens: tokens,
	}
	repository.AddUser(u)
	log.Println("Len(AuthenticatedUsers) : " + strconv.Itoa(len(repository.GetUsers())) + "\n")
	return sessionCode
}

// Removes the user session based on his session_code
func RemoveSessionCode(sessionCode string) {
	var temp []*models.User
	for _, u := range repository.GetUsers() {
		if u.Code != sessionCode {
			temp = append(temp, u)
		}
	}
	repository.ReplaceAllUsers(temp)
	log.Println("Len(AuthenticatedUsers) : " + strconv.Itoa(len(repository.GetUsers())) + "\n")
}

// Print the Auth Middleware messages
func logInAuthMiddleware(msg string) {
	log.Println(prefixAuthMiddleware + " " + msg + "\n")
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
