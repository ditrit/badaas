package router

import (
	"net/http"

	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter(
	//middlewares
	jsonController middlewares.JSONController,
	middlewareLogger middlewares.MiddlewareLogger,
	authenticationMiddleware middlewares.AuthenticationMiddleware,

	// controllers
	basicAuthenticationController controllers.BasicAuthenticationController,
	informationController controllers.InformationController,
	eavController controllers.EAVController,
) http.Handler {
	router := mux.NewRouter()
	router.Use(middlewareLogger.Handle)

	router.HandleFunc(
		"/info",
		jsonController.Wrap(informationController.Info),
	).Methods("GET")
	router.HandleFunc(
		"/login",
		jsonController.Wrap(
			basicAuthenticationController.BasicLoginHandler,
		),
	).Methods("POST")

	protected := router.PathPrefix("").Subrouter()
	protected.Use(authenticationMiddleware.Handle)

	protected.HandleFunc("/logout", jsonController.Wrap(basicAuthenticationController.Logout)).Methods("GET")

	// Objects CRUD
	objectsBase := "/objects/{type}"
	objectsWithID := objectsBase + "/{id}"
	router.HandleFunc(objectsWithID, jsonController.Wrap(eavController.GetObject)).Methods("GET")
	router.HandleFunc(objectsBase, jsonController.Wrap(eavController.GetObjects)).Methods("GET")
	router.HandleFunc(objectsBase, jsonController.Wrap(eavController.CreateObject)).Methods("POST")
	router.HandleFunc(objectsWithID, jsonController.Wrap(eavController.UpdateObject)).Methods("PUT")
	router.HandleFunc(objectsWithID, jsonController.Wrap(eavController.DeleteObject)).Methods("DELETE")

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
			"Access-Control-Request-Headers", "Access-Control-Request-Method",
			"Connection", "Host", "Origin", "User-Agent", "Referer",
			"Cache-Control", "X-header",
		}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(0),
	)(router)

	return cors
}
