package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/controllers"
	mockServices "github.com/ditrit/badaas/mocks/services"
	"github.com/ditrit/badaas/persistence/models"
)

var logger, _ = zap.NewDevelopment()

// ----------------------- GetObject -----------------------

func TestGetWithoutTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/",
		strings.NewReader(""),
	)

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	uuid := badorm.NewUUID()
	eavService.
		On("GetEntity", "no-exists", uuid).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/no-exists/id",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists", "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithoutEntityIDReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithEntityIDNotUUIDReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/not-uuid",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": "not-uuid"})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrIDNotAnUUID)
}

func TestGetWithEntityIDThatDoesNotExistReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()

	eavService.
		On("GetEntity", entityType.Name, uuid).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithErrorInDBReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()

	eavService.
		On("GetEntity", entityType.Name, uuid).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestGetWithCorrectIDReturnsObject(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity := &models.Entity{
		EntityType: entityType,
	}

	uuid := badorm.NewUUID()

	eavService.
		On("GetEntity", entityType.Name, uuid).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	entityReturned, err := controller.GetObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, entityReturned)
}

// ----------------------- GetAll -----------------------

func TestGetAllWithoutTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/",
		strings.NewReader(""),
	)

	_, err := controller.GetObjects(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetAllOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	eavService.
		On("GetEntities", "no-exists", map[string]any{}).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/no-exists",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists"})

	_, err := controller.GetObjects(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetAllWithErrorInDBReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	eavService.
		On("GetEntities", "no-exists", map[string]any{}).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/no-exists",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists"})

	_, err := controller.GetObjects(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestGetAllWithoutParams(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

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
		On("GetEntities", entityType.Name, map[string]any{}).
		Return([]*models.Entity{entity1, entity2}, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	entitiesReturned, err := controller.GetObjects(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 2)
	assert.Contains(t, entitiesReturned, entity1)
	assert.Contains(t, entitiesReturned, entity2)
}

func TestGetAllWithParams(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity1 := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("GetEntities", entityType.Name, map[string]any{"param1": "something"}).
		Return([]*models.Entity{entity1}, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/objects/exists/",
		strings.NewReader("{\"param1\": \"something\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	entitiesReturned, err := controller.GetObjects(response, request)
	assert.Nil(t, err)
	assert.Len(t, entitiesReturned, 1)
	assert.Contains(t, entitiesReturned, entity1)
}

// ----------------------- DeleteObject -----------------------

func TestDeleteWithoutTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodDelete,
		"/objects/",
		strings.NewReader(""),
	)

	_, err := controller.DeleteObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestDeleteOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	uuid := badorm.NewUUID()

	eavService.
		On("DeleteEntity", "no-exists", uuid).
		Return(gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodDelete,
		"/objects/no-exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists", "id": uuid.String()})

	_, err := controller.DeleteObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestDeleteObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()

	eavService.
		On("DeleteEntity", entityType.Name, uuid).
		Return(errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodDelete,
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.DeleteObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestDeleteObjectReturnsNil(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()

	eavService.
		On("DeleteEntity", entityType.Name, uuid).
		Return(nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodDelete,
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	returned, err := controller.DeleteObject(response, request)
	assert.Nil(t, err)
	assert.Nil(t, returned)
}

// ----------------------- CreateObject -----------------------

func TestCreateWithoutTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/objects/",
		strings.NewReader(""),
	)

	_, err := controller.CreateObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestCreateObjectWithBadJSONReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/objects/exists",
		strings.NewReader("bad json"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	_, err := controller.CreateObject(response, request)
	assert.ErrorContains(t, err, "The schema of the received data is not correct")
}

func TestCreateOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	eavService.
		On("CreateEntity", "no-exists", map[string]any{"1": "1"}).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/objects/no-exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists"})

	_, err := controller.CreateObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestCreteObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	eavService.
		On("CreateEntity", entityType.Name, map[string]any{"1": "1"}).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/objects/exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	_, err := controller.CreateObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestCreteObjectReturnsObject(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	entity := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("CreateEntity", entityType.Name, map[string]any{"1": "1"}).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/objects/exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	responded, err := controller.CreateObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, responded)
}

// ----------------------- UpdateObject -----------------------

func TestModifyWithoutTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPut,
		"/objects/",
		strings.NewReader(""),
	)

	_, err := controller.UpdateObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestUpdateObjectWithBadJSONReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPut,
		"/objects/exists",
		strings.NewReader("bad json"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.UpdateObject(response, request)
	assert.ErrorContains(t, err, "The schema of the received data is not correct")
}

func TestModifyOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	uuid := badorm.NewUUID()

	eavService.
		On("UpdateEntity", "no-exists", uuid, map[string]any{"1": "1"}).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/objects/no-exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists", "id": uuid.String()})

	_, err := controller.UpdateObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestUpdateObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()

	eavService.
		On("UpdateEntity", entityType.Name, uuid, map[string]any{"1": "1"}).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPut,
		"/objects/exists/"+uuid.String(),
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.UpdateObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestUpdateObjectReturnsObject(t *testing.T) {
	eavService := mockServices.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := badorm.NewUUID()
	entity := &models.Entity{
		EntityType: entityType,
	}

	eavService.
		On("UpdateEntity", entityType.Name, uuid, map[string]any{"1": "1"}).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPut,
		"/objects/exists/"+uuid.String(),
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	responded, err := controller.UpdateObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, responded)
}
