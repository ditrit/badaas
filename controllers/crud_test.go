package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/controllers"
	mockBadorm "github.com/ditrit/badaas/mocks/badorm"
	mockUnsafe "github.com/ditrit/badaas/mocks/badorm/unsafe"
)

// ----------------------- GetModel -----------------------

type Model struct {
	badorm.UUIDModel
}

func TestCRUDGetWithoutEntityIDReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)

	_, err := route.Controller.GetModel(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestCRUDGetWithEntityIDNotUUIDReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"id": "not-uuid"})

	_, err := route.Controller.GetModel(response, request)
	assert.ErrorIs(t, err, controllers.ErrIDNotAnUUID)
}

func TestCRUDGetWithEntityIDThatDoesNotExistReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	uuid := badorm.NewUUID()

	crudService.
		On("GetByID", uuid).
		Return(nil, gorm.ErrRecordNotFound)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{"id": uuid.String()})

	_, err := route.Controller.GetModel(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestCRUDGetWithErrorInDBReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	uuid := badorm.NewUUID()

	crudService.
		On("GetByID", uuid).
		Return(nil, errors.New("db error"))

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{"id": uuid.String()})

	_, err := route.Controller.GetModel(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestCRUDGetWithCorrectIDReturnsObject(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	uuid := badorm.NewUUID()
	entity := Model{}

	crudService.
		On("GetByID", uuid).
		Return(&entity, nil)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{"id": uuid.String()})
	entityReturned, err := route.Controller.GetModel(response, request)
	assert.Nil(t, err)
	assert.Equal(t, &entity, entityReturned)
}

// ----------------------- GetModels -----------------------

func TestGetModelsWithErrorInDBReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	crudUnsafeService.
		On("Query", map[string]any{}).
		Return(nil, errors.New("db error"))

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)

	_, err := route.Controller.GetModels(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestGetModelsWithoutParams(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	entity1 := &Model{}
	entity2 := &Model{}

	crudUnsafeService.
		On("Query", map[string]any{}).
		Return([]*Model{entity1, entity2}, nil)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)

	entitiesReturned, err := route.Controller.GetModels(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 2)
	assert.Contains(t, entitiesReturned, entity1)
	assert.Contains(t, entitiesReturned, entity2)
}

func TestGetModelsWithParams(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	entity1 := &Model{}

	crudUnsafeService.
		On("Query", map[string]any{"param1": "something"}).
		Return([]*Model{entity1}, nil)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader("{\"param1\": \"something\"}"),
	)

	entitiesReturned, err := route.Controller.GetModels(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 1)
	assert.Contains(t, entitiesReturned, entity1)
}

func TestGetModelsWithParamsNotJsonReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, badorm.UUID](t)
	crudUnsafeService := mockUnsafe.NewCRUDService[Model, badorm.UUID](t)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
		crudUnsafeService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader("bad json"),
	)

	_, err := route.Controller.GetModels(response, request)
	assert.ErrorIs(t, err, controllers.HTTPErrRequestMalformed)
}
