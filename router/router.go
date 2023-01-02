package router

import (
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter(
	// configuration holders
	authenticationConfiguration configuration.AuthenticationConfiguration,

	//middlewares
	jsonController middlewares.JSONController,
	middlewareLogger middlewares.MiddlewareLogger,
	authenticationMiddleware middlewares.AuthenticationMiddleware,

	// controllers
	basicAuthentificationController controllers.BasicAuthentificationController,
	informationController controllers.InformationController,
	oidcController controllers.OIDCController,
) http.Handler {
	router := mux.NewRouter() //.PathPrefix(fmt.Sprintf("/%v", resources.Version)).Subrouter()
	router.Use(middlewareLogger.Handle)

	router.HandleFunc(
		"/info",
		jsonController.Wrap(informationController.Info),
	).Methods("GET")
	router.HandleFunc(
		"/auth/basic/login",
		jsonController.Wrap(
			basicAuthentificationController.BasicLoginHandler,
		),
	).Methods("POST")

	// OIDC
	if authenticationConfiguration.GetAuthType() == configuration.AuthTypeOIDC {
		router.HandleFunc(
			"/auth/oidc/redirect-url",
			jsonController.Wrap(oidcController.RedirectURL),
		).Methods("GET")
		router.HandleFunc(
			"/auth/oidc/callback",
			jsonController.Wrap(oidcController.CallBack),
		).Methods("GET")
	}

	protected := router.PathPrefix("").Subrouter()
	protected.Use(authenticationMiddleware.Handle)

	protected.HandleFunc(
		"/auth/logout",
		jsonController.Wrap(basicAuthentificationController.Logout),
	).Methods("GET")

	return router
}
