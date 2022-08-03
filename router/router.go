package router

import (
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/controllers/openid_connect"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/info", controllers.Info).Methods("GET")

	oidcRouter := router.PathPrefix("/oidc").Subrouter()
	createOidcSubrouter(oidcRouter)

	return router
}

func createOidcSubrouter(subRouter *mux.Router) {

	subRouter.Use(openid_connect.MiddlewareLogger)
	subRouter.HandleFunc("/login-screen", controllers.LoginScreen).Methods("GET")
	subRouter.HandleFunc("/get-session-code", controllers.GetSessionCode).Methods("POST")
	subRouter.HandleFunc("/refresh-tokens", controllers.RefreshTokens).Methods("GET")

	// Protected routes : a valid session_code is needed to enter these routes
	protectedRouter := subRouter.PathPrefix("/").Subrouter()
	protectedRouter.Use(openid_connect.MiddlewareAuthenticator)
	protectedRouter.HandleFunc("/authenticated", controllers.Authenticated).Methods("GET")
	protectedRouter.HandleFunc("/logout", controllers.Logout).Methods("GET")

}
