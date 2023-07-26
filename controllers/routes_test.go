package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	mocks "github.com/ditrit/badaas/mocks/configuration"
	mockControllers "github.com/ditrit/badaas/mocks/controllers"
	mockMiddlewares "github.com/ditrit/badaas/mocks/router/middlewares"
	mockUserServices "github.com/ditrit/badaas/mocks/services/userservice"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/router/middlewares"
)

func TestCreateSuperUser(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	initializationConfig := mocks.NewInitializationConfiguration(t)
	initializationConfig.On("GetAdminPassword").Return("adminpassword")
	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, nil)
	err := createSuperUser(
		logger,
		initializationConfig,
		userService,
	)
	assert.NoError(t, err)
}

func TestCreateSuperUser_UserExists(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	initializationConfig := mocks.NewInitializationConfiguration(t)
	initializationConfig.On("GetAdminPassword").Return("adminpassword")
	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, errors.New("user already exist in database"))
	err := createSuperUser(
		logger,
		initializationConfig,
		userService,
	)
	assert.NoError(t, err)

	require.Equal(t, 1, logs.Len())
}

func TestCreateSuperUser_UserServiceError(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	initializationConfig := mocks.NewInitializationConfiguration(t)
	initializationConfig.On("GetAdminPassword").Return("adminpassword")
	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, errors.New("email not valid"))
	err := createSuperUser(
		logger,
		initializationConfig,
		userService,
	)
	assert.Error(t, err)

	require.Equal(t, 1, logs.Len())
}

var logger, _ = zap.NewDevelopment()

func TestAddInfoRoutes(t *testing.T) {
	jsonController := middlewares.NewJSONController(logger)
	informationController := NewInfoController(semver.MustParse("1.0.1"))

	router := router.NewRouter()
	AddInfoRoutes(
		router,
		jsonController,
		informationController,
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/info",
		nil,
	)

	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), "{\"status\":\"OK\",\"version\":\"1.0.1\"}")
}

func TestAddLoginRoutes(t *testing.T) {
	jsonController := middlewares.NewJSONController(logger)

	initializationConfig := mocks.NewInitializationConfiguration(t)
	initializationConfig.
		On("GetAdminPassword").Return("adminpassword")

	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, nil)

	basicAuthenticationController := mockControllers.NewBasicAuthenticationController(t)
	basicAuthenticationController.
		On("BasicLoginHandler", mock.Anything, mock.Anything).
		Return(map[string]string{"login": "called"}, nil)

	authenticationMiddleware := mockMiddlewares.NewAuthenticationMiddleware(t)

	router := router.NewRouter()
	AddAuthRoutes(
		nil,
		router,
		authenticationMiddleware,
		basicAuthenticationController,
		jsonController,
		initializationConfig,
		userService,
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/login",
		nil,
	)

	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), "{\"login\":\"called\"}")
}
