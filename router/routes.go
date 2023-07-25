package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/ditrit/badaas/services/userservice"
)

func AddInfoRoutes(
	router *mux.Router,
	jsonController middlewares.JSONController,
	infoController controllers.InformationController,
) {
	router.HandleFunc(
		"/info",
		jsonController.Wrap(infoController.Info),
	).Methods(http.MethodGet)
}

// Adds to the "router" the routes for handling authentication:
// /login
// /logout
// And creates a very first user
func AddAuthRoutes(
	logger *zap.Logger,
	router *mux.Router,
	authenticationMiddleware middlewares.AuthenticationMiddleware,
	basicAuthenticationController controllers.BasicAuthenticationController,
	jsonController middlewares.JSONController,
	config configuration.InitializationConfiguration,
	userService userservice.UserService,
) error {
	router.HandleFunc(
		"/login",
		jsonController.Wrap(basicAuthenticationController.BasicLoginHandler),
	).Methods(http.MethodPost)

	protected := router.PathPrefix("").Subrouter()
	protected.Use(authenticationMiddleware.Handle)

	protected.HandleFunc(
		"/logout",
		jsonController.Wrap(basicAuthenticationController.Logout),
	).Methods(http.MethodGet)

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
