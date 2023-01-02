package router

import (
	"net/http"
	"testing"

	"github.com/ditrit/badaas/configuration"
	configurationMocks "github.com/ditrit/badaas/mocks/configuration"
	controllersMocks "github.com/ditrit/badaas/mocks/controllers"
	middlewaresMocks "github.com/ditrit/badaas/mocks/router/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetupRouter(t *testing.T) {
	jsonController := middlewaresMocks.NewJSONController(t)
	middlewareLogger := middlewaresMocks.NewMiddlewareLogger(t)
	authenticationMiddleware := middlewaresMocks.NewAuthenticationMiddleware(t)

	authenticationConfig := configurationMocks.NewAuthenticationConfiguration(t)
	authenticationConfig.On("GetAuthType").Return(configuration.AuthTypeOIDC)

	basicController := controllersMocks.NewBasicAuthentificationController(t)
	informationController := controllersMocks.NewInformationController(t)
	jsonController.On("Wrap", mock.Anything).Return(func(response http.ResponseWriter, request *http.Request) {})
	oidcController := controllersMocks.NewOIDCController(t)
	router := SetupRouter(authenticationConfig, jsonController, middlewareLogger, authenticationMiddleware, basicController, informationController, oidcController)
	assert.NotNil(t, router)
}
