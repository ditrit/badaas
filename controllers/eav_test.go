package controllers_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ditrit/badaas/controllers"
	mocksEAVService "github.com/ditrit/badaas/mocks/services/eavservice"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/eavservice"
	uuid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var logger, _ = zap.NewDevelopment()

// ----------------------- GetObject -----------------------

func TestGetWithoutTypeReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/",
		strings.NewReader(""),
	)

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	eavService.
		On("GetEntityTypeByName", "no-exists").
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/no-exists/id",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists"})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetWithoutEntityIDReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists"})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithEntityIDNotUUIDReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/not-uuid",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": "not-uuid"})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrIDNotAnUUID)
}

func TestGetWithEntityIDThatDoesNotExistReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithEntityIDThatDoesNotMatchEntityTypeReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(nil, eavservice.ErrIDDontMatchEntityType)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeDontMatchEntity)
}

func TestGetWithErrorInDBReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestGetWithCorrectIDReturnsObject(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity := &models.Entity{
		EntityType: entityType,
	}

	uuid := uuid.New()

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	entityReturned, err := controller.GetObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, entityReturned)
}

// ----------------------- GetAll -----------------------

func TestGetAllOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	eavService.
		On("GetEntityTypeByName", "no-exists").
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/no-exists",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists"})

	_, err := controller.GetAll(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetAllWithoutParams(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity1 := &models.Entity{
		EntityType: entityType,
	}
	entity2 := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntitiesWithParams", entityType, map[string]string{}).
		Return([]*models.Entity{entity1, entity2})

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists"})

	entitiesReturned, err := controller.GetAll(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 2)
	assert.Contains(t, entitiesReturned, entity1)
	assert.Contains(t, entitiesReturned, entity2)
}

func TestGetAllWithParams(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity1 := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntitiesWithParams", entityType, map[string]string{"param1": "something"}).
		Return([]*models.Entity{entity1})

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/v1/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists"})
	q := request.URL.Query()
	q.Add("param1", "something")
	request.URL.RawQuery = q.Encode()

	entitiesReturned, err := controller.GetAll(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 1)
	assert.Contains(t, entitiesReturned, entity1)
}

// ----------------------- DeleteObject -----------------------

func TestDeleteObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity := &models.Entity{
		EntityType: entityType,
	}

	uuid := uuid.New()

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(entity, nil)

	eavService.
		On("DeleteEntity", entity).
		Return(errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"DELETE",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	_, err := controller.DeleteObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestDeleteObjectReturnsNil(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity := &models.Entity{
		EntityType: entityType,
	}

	uuid := uuid.New()

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(entity, nil)

	eavService.
		On("DeleteEntity", entity).
		Return(nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"DELETE",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	returned, err := controller.DeleteObject(response, request)
	assert.Nil(t, err)
	assert.Nil(t, returned)
}

// ----------------------- CreateObject -----------------------

func TestCreateObjectWithBadJSONReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/objects/exists",
		strings.NewReader("bad json"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists"})

	_, err := controller.CreateObject(response, request)
	assert.ErrorContains(t, err, "json decoding failed")
}

func TestCreteObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("CreateEntity", entityType, map[string]any{"1": "1"}).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/objects/exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists"})

	_, err := controller.CreateObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestCreteObjectReturnsObject(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("CreateEntity", entityType, map[string]any{"1": "1"}).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/v1/objects/exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists"})

	responded, err := controller.CreateObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, responded)
}

// ----------------------- ModifyObject -----------------------

func TestModifyObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()
	entity := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(entity, nil)

	eavService.
		On("UpdateEntity", entity, map[string]any{"1": "1"}).
		Return(errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"PUT",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	_, err := controller.ModifyObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestModifyObjectReturnsObject(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()
	entity := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("GetEntityTypeByName", "exists").
		Return(entityType, nil)

	eavService.
		On("GetEntity", entityType, uuid).
		Return(entity, nil)

	eavService.
		On("UpdateEntity", entity, map[string]any{"1": "1"}).
		Return(nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"PUT",
		"/v1/objects/exists/"+uuid.String(),
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "exists", "id": uuid.String()})

	responded, err := controller.ModifyObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, responded)
}
