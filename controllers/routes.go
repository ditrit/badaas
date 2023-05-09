package controllers

import (
	"strings"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func AddInfoRoutes(
	router *mux.Router,
	jsonController middlewares.JSONController,
	infoController InformationController,
) {
	router.HandleFunc(
		"/info",
		jsonController.Wrap(infoController.Info),
	).Methods("GET")
}

// Adds to the "router" the routes for handling authentication:
// /login
// /logout
// And creates a very first user
func AddAuthRoutes(
	logger *zap.Logger,
	router *mux.Router,
	authenticationMiddleware middlewares.AuthenticationMiddleware,
	basicAuthenticationController BasicAuthenticationController,
	jsonController middlewares.JSONController,
	config configuration.InitializationConfiguration,
	userService userservice.UserService,
) error {
	router.HandleFunc(
		"/login",
		jsonController.Wrap(basicAuthenticationController.BasicLoginHandler),
	).Methods("POST")

	protected := router.PathPrefix("").Subrouter()
	protected.Use(authenticationMiddleware.Handle)

	protected.HandleFunc(
		"/logout",
		jsonController.Wrap(basicAuthenticationController.Logout),
	).Methods("GET")

	return createSuperUser(logger, config, userService)
}

// Create a super user
func createSuperUser(
	logger *zap.Logger,
	config configuration.InitializationConfiguration,
	userService userservice.UserService,
) error {
	_, err := userService.NewUser("admin", "admin-no-reply@badaas.com", config.GetAdminPassword())
	if err != nil {
		if !strings.Contains(err.Error(), "already exist in database") {
			logger.Sugar().Errorf("failed to save the super admin %w", err)
			return err
		}
		logger.Sugar().Infof("The superadmin user already exists in database")
	}

	return nil
}

func AddEAVCRUDRoutes(
	router *mux.Router,
	eavController CRUDController,
	jsonController middlewares.JSONController,
) {
	// Objects CRUD
	objectsBase := "/eav/objects/{type}"
	objectsWithID := objectsBase + "/{id}"
	router.HandleFunc(objectsWithID, jsonController.Wrap(eavController.GetObject)).Methods("GET")
	router.HandleFunc(objectsBase, jsonController.Wrap(eavController.GetObjects)).Methods("GET")
	router.HandleFunc(objectsBase, jsonController.Wrap(eavController.CreateObject)).Methods("POST")
	router.HandleFunc(objectsWithID, jsonController.Wrap(eavController.UpdateObject)).Methods("PUT")
	router.HandleFunc(objectsWithID, jsonController.Wrap(eavController.DeleteObject)).Methods("DELETE")
}

func AddCRUDRoutes(
	router *mux.Router,
	generalCRUDController *GeneralCRUDController,
	jsonController middlewares.JSONController,
) {
	// Objects CRUD
	objectsBase := "/objects/{type}"
	objectsWithID := objectsBase + "/{id}"
	router.HandleFunc(objectsWithID, jsonController.Wrap(generalCRUDController.GetObject)).Methods("GET")
	router.HandleFunc(objectsBase, jsonController.Wrap(generalCRUDController.GetObjects)).Methods("GET")
	router.HandleFunc(objectsBase, jsonController.Wrap(generalCRUDController.CreateObject)).Methods("POST")
	router.HandleFunc(objectsWithID, jsonController.Wrap(generalCRUDController.UpdateObject)).Methods("PUT")
	router.HandleFunc(objectsWithID, jsonController.Wrap(generalCRUDController.DeleteObject)).Methods("DELETE")
}
