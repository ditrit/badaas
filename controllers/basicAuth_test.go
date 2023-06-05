package controllers_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/httperrors"
	mocksSessionService "github.com/ditrit/badaas/mocks/services/sessionservice"
	mocksUserService "github.com/ditrit/badaas/mocks/services/userservice"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

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
		"POST",
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
		Return(nil, httperrors.AnError)
	sessionService := mocksSessionService.NewSessionService(t)

	controller := controllers.NewBasicAuthenticationController(
		logger,
		userService,
		sessionService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
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
		"POST",
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
		Return(nil, httperrors.AnError)

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
		"POST",
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
	assert.Equal(t, payload, dto.DTOLoginSuccess{
		Email:    "bob@email.com",
		ID:       user.ID.String(),
		Username: user.Username,
	})
}
