package router

import (
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.LoggerMW)

	router.HandleFunc("/info", controllers.Info).Methods("GET")

	// basic auth login
	router.HandleFunc("/login", controllers.BasicLoginHandler).Methods("POST")

	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middlewares.AuthenticationMW)
	// JWT refresh
	protectedRouter.HandleFunc("/refresh", controllers.RefreshJWT).Methods("GET")

	return router
}
