package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Creates the backend server with the necessary endpoints to implement the OIDC protocol with Google and Gitlab as providers
func main() {
	// Unprotected routes : no session-code needed to enter these routes
	router := mux.NewRouter()
	router.Use(MiddlewareLogger)
	router.HandleFunc("/login-screen", LoginScreen).Methods("GET")
	router.HandleFunc("/get-session-code", GetSessionCode).Methods("POST")
	router.HandleFunc("/refresh-tokens", RefreshTokens).Methods("GET")

	// Protected routes : a session-code is needed to enter these routes
	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(MiddlewareAuthenticator)
	protectedRouter.HandleFunc("/authenticated", Authenticated).Methods("GET")
	protectedRouter.HandleFunc("/logout", Logout).Methods("GET")

	// It may be a good idea to choose the CORS options at the bare minimum level
	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(0),
	)(router)
	fmt.Println("Ready !")
	http.ListenAndServe(":8090", cors)
}
