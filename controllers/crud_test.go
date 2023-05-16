package controllers_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ditrit/badaas/controllers"
	mockBadorm "github.com/ditrit/badaas/mocks/badorm"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ----------------------- GetObject -----------------------

type Model struct {
	ID uuid.UUID
}

func TestCRUDGetWithoutEntityIDReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{})

	_, err := route.Controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestCRUDGetWithEntityIDNotUUIDReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"id": "not-uuid"})

	_, err := route.Controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrIDNotAnUUID)
}

func TestCRUDGetWithEntityIDThatDoesNotExistReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	uuid := uuid.New()

	crudService.
		On("GetEntity", uuid).
		Return(nil, gorm.ErrRecordNotFound)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{"id": uuid.String()})

	_, err := route.Controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestCRUDGetWithErrorInDBReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	uuid := uuid.New()

	crudService.
		On("GetEntity", uuid).
		Return(nil, errors.New("db error"))

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{"id": uuid.String()})

	_, err := route.Controller.GetObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestCRUDGetWithCorrectIDReturnsObject(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	uuid := uuid.New()
	entity := Model{}

	crudService.
		On("GetEntity", uuid).
		Return(&entity, nil)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{"id": uuid.String()})
	entityReturned, err := route.Controller.GetObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, &entity, entityReturned)
}

// ----------------------- GetEntities -----------------------

func TestGetEntitiesWithErrorInDBReturnsError(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	crudService.
		On("GetEntities", map[string]any{}).
		Return(nil, errors.New("db error"))

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{})
	_, err := route.Controller.GetObjects(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestGetEntitiesWithoutParams(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	entity1 := &Model{}
	entity2 := &Model{}

	crudService.
		On("GetEntities", map[string]any{}).
		Return([]*Model{entity1, entity2}, nil)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)

	request = mux.SetURLVars(request, map[string]string{})
	entitiesReturned, err := route.Controller.GetObjects(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 2)
	assert.Contains(t, entitiesReturned, entity1)
	assert.Contains(t, entitiesReturned, entity2)
}

func TestGetEntitiesWithParams(t *testing.T) {
	crudService := mockBadorm.NewCRUDService[Model, uuid.UUID](t)

	entity1 := &Model{}

	crudService.
		On("GetEntities", map[string]any{"param1": "something"}).
		Return([]*Model{entity1}, nil)

	route := controllers.NewCRUDController[Model](
		logger,
		crudService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader("{\"param1\": \"something\"}"),
	)

	request = mux.SetURLVars(request, map[string]string{})
	entitiesReturned, err := route.Controller.GetObjects(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 1)
	assert.Contains(t, entitiesReturned, entity1)
}
