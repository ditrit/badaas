package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ditrit/badaas/httperrors"
	mockRepository "github.com/ditrit/badaas/mocks/persistence/repository"
	mockOIDCService "github.com/ditrit/badaas/mocks/services/auth/protocols/oidcservice"
	mockSessionService "github.com/ditrit/badaas/mocks/services/sessionservice"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewOIDCController(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)
	assert.NotNil(t, oidcController)
}

func TestOIDCController_RedirectURL(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/redirect-url",
		strings.NewReader(""),
	)

	oidcService.
		On("BuildRedirectURL", mock.AnythingOfType("string")).
		Return("https://account.domain.tld")

	payload, err := oidcController.RedirectURL(response, request)
	require.NoError(t, err)
	require.NotNil(t, payload)
	cookies := response.Result().Cookies()

	assert.Len(t, cookies, 1)
	assert.NotNil(t, cookies[0])
}

func TestOIDCController_CallBack_noState(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback",
		strings.NewReader(""),
	)

	payload, err := oidcController.CallBack(response, request)
	require.Equal(t, err, ErrMissingCookieState)
	require.Nil(t, payload)
}

func TestOIDCController_CallBack_stateDontMatch(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback",
		strings.NewReader(""),
	)

	// oidcService.
	// On("ExchangeAuthorizationCode", mock.AnythingOfType("string")).
	// Return("user@dev.com", nil)
	request.AddCookie(createCallbackCookie("state", "6543", request))

	payload, err := oidcController.CallBack(response, request)
	require.Equal(t, ErrStateDidNotMatch, err)
	require.Nil(t, payload)
}

func TestOIDCController_CallBack_noCode(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback?state=6543",
		strings.NewReader(""),
	)

	// oidcService.
	// On("ExchangeAuthorizationCode", mock.AnythingOfType("string")).
	// Return("user@dev.com", nil)
	request.AddCookie(createCallbackCookie("state", "6543", request))

	payload, err := oidcController.CallBack(response, request)
	require.Equal(t, ErrAuthorizationCodeNotFound, err)
	require.Nil(t, payload)
}

func TestOIDCController_CallBack_exchangeFail(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback?state=6543&code=abcd",
		strings.NewReader(""),
	)

	oidcService.
		On("ExchangeAuthorizationCode", "abcd").
		Return("", httperrors.AnError)
	request.AddCookie(createCallbackCookie("state", "6543", request))

	payload, err := oidcController.CallBack(response, request)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestOIDCController_CallBack_userNotFound(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback?state=6543&code=abcd",
		strings.NewReader(""),
	)

	oidcService.
		On("ExchangeAuthorizationCode", "abcd").
		Return("user@dev.com", nil)
	userRepository.
		On("Find", mock.Anything, nil, nil).
		Return(pagination.NewPage([]*models.User{}, 0, 10, 0), nil)
	request.AddCookie(createCallbackCookie("state", "6543", request))

	payload, err := oidcController.CallBack(response, request)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestOIDCController_CallBack_userError(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback?state=6543&code=abcd",
		strings.NewReader(""),
	)

	oidcService.
		On("ExchangeAuthorizationCode", "abcd").
		Return("user@dev.com", nil)
	userRepository.
		On("Find", mock.Anything, nil, nil).
		Return(nil, httperrors.AnError)
	request.AddCookie(createCallbackCookie("state", "6543", request))

	payload, err := oidcController.CallBack(response, request)
	require.Equal(t, httperrors.AnError, err)
	require.Nil(t, payload)
}

func TestOIDCController_CallBack_success(t *testing.T) {
	oidcService := mockOIDCService.NewOIDCService(t)
	sessionService := mockSessionService.NewSessionService(t)
	userRepository := mockRepository.NewCRUDRepository[models.User, uint](t)
	oidcController := NewOIDCController(zap.L(), oidcService, userRepository, sessionService)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/callback?state=6543&code=abcd",
		strings.NewReader(""),
	)
	user := &models.User{Username: "bob"}
	oidcService.
		On("ExchangeAuthorizationCode", "abcd").
		Return("user@dev.com", nil)
	userRepository.
		On("Find", mock.Anything, nil, nil).
		Return(pagination.NewPage([]*models.User{user},
			0, 10, 0), nil)
	request.AddCookie(createCallbackCookie("state", "6543", request))
	sessionService.On("LogUserIn", user, response).Return(nil)
	payload, err := oidcController.CallBack(response, request)
	require.NoError(t, err)
	require.Nil(t, payload)
}

func Test_randString(t *testing.T) {
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
	assert.NotEqual(t, randString(10), randString(10))
}

func Test_setCallbackCookie(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/auth/oidc/redirect-url",
		strings.NewReader(""),
	)
	setCallbackCookie(response, request, "state", "6543")
	cookies := response.Result().Cookies()
	require.Len(t, cookies, 1)
	c := cookies[0]
	require.NotNil(t, c)
	assert.Equal(t, "state", c.Name)
	assert.Equal(t, "6543", c.Value)
	assert.Equal(t, "/", c.Path)
	assert.Equal(t,
		int(time.Duration(15*time.Minute).Seconds()),
		c.MaxAge,
	)
	assert.Equal(t, true, c.HttpOnly)
	assert.Equal(t, false, c.Secure)
}
