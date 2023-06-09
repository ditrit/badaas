package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/httperrors"
	mocksSessionService "github.com/ditrit/badaas/mocks/services/sessionservice"
	mocksUserService "github.com/ditrit/badaas/mocks/services/userservice"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
)

// ExampleErr is an HTTPError instance useful for testing.  If the code does not care
// about HTTPError specifics, and only needs to return the HTTPError for example, this
// HTTPError should be used to make the test code more readable.
var ExampleErr httperrors.HTTPError = &httperrors.HTTPErrorImpl{
	Status:      -1,
	Err:         "TESTING ERROR",
	Message:     "USE ONLY FOR TESTING",
	GolangError: nil,
}

func Test_BasicLoginHandler_MalformedRequest(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	userService := mocksUserService.NewUserService(t)
	sessionService := mocksSessionService.NewSessionService(t)

	controller := controllers.NewBasicAuthenticationController(
		logger,
		userService,
		sessionService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader("qsdqsdqsd"),
	)

	payload, err := controller.BasicLoginHandler(response, request)
	assert.Equal(t, controllers.HTTPErrRequestMalformed, err)
	assert.Nil(t, payload)
}

func Test_BasicLoginHandler_UserNotFound(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	loginJSONStruct := dto.UserLoginDTO{
		Email:    "bob@email.com",
		Password: "1234",
	}
	userService := mocksUserService.NewUserService(t)
	userService.
		On("GetUser", loginJSONStruct).
		Return(nil, ExampleErr)

	sessionService := mocksSessionService.NewSessionService(t)

	controller := controllers.NewBasicAuthenticationController(
		logger,
		userService,
		sessionService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
			"email": "bob@email.com",
			"password":"1234"
		}`),
	)

	payload, err := controller.BasicLoginHandler(response, request)
	assert.Error(t, err)
	assert.Nil(t, payload)
}

func Test_BasicLoginHandler_LoginFailed(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	loginJSONStruct := dto.UserLoginDTO{
		Email:    "bob@email.com",
		Password: "1234",
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
			"email": "bob@email.com",
			"password":"1234"
		}`),
	)
	userService := mocksUserService.NewUserService(t)
	user := &models.User{
		UUIDModel: badorm.UUIDModel{},
		Username:  "bob",
		Email:     "bob@email.com",
		Password:  []byte("hash of 1234"),
	}
	userService.
		On("GetUser", loginJSONStruct).
		Return(user, nil)

	sessionService := mocksSessionService.NewSessionService(t)
	sessionService.
		On("LogUserIn", user).
		Return(nil, ExampleErr)

	controller := controllers.NewBasicAuthenticationController(
		logger,
		userService,
		sessionService,
	)

	payload, err := controller.BasicLoginHandler(response, request)
	assert.Error(t, err)
	assert.Nil(t, payload)
}

func Test_BasicLoginHandler_LoginSuccess(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	loginJSONStruct := dto.UserLoginDTO{
		Email:    "bob@email.com",
		Password: "1234",
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(`{
			"email": "bob@email.com",
			"password":"1234"
		}`),
	)
	userService := mocksUserService.NewUserService(t)
	user := &models.User{
		UUIDModel: badorm.UUIDModel{
			ID: badorm.NilUUID,
		},
		Username: "bob",
		Email:    "bob@email.com",
		Password: []byte("hash of 1234"),
	}
	userService.
		On("GetUser", loginJSONStruct).
		Return(user, nil)

	sessionService := mocksSessionService.NewSessionService(t)
	sessionService.
		On("LogUserIn", user).
		Return(models.NewSession(user.ID, time.Duration(5)), nil)

	controller := controllers.NewBasicAuthenticationController(
		logger,
		userService,
		sessionService,
	)

	payload, err := controller.BasicLoginHandler(response, request)
	assert.NoError(t, err)
	assert.Equal(t, payload, dto.LoginSuccess{
		Email:    "bob@email.com",
		ID:       user.ID.String(),
		Username: user.Username,
	})
}
