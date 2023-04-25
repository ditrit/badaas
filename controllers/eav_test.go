package controllers_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ditrit/badaas/controllers"
	mocksEAVService "github.com/ditrit/badaas/mocks/services/eavservice"
	"github.com/ditrit/badaas/persistence/models"
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
		"/objects/",
		strings.NewReader(""),
	)

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityTypeNotFound)
}

func TestGetOfNotExistentTypeReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	uuid := uuid.New()
	eavService.
		On("GetEntity", "no-exists", uuid).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/no-exists/id",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists", "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithoutEntityIDReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithEntityIDNotUUIDReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/not-uuid",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": "not-uuid"})

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
		On("GetEntity", entityType.Name, uuid).
		Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.GetObject(response, request)
	assert.ErrorIs(t, err, controllers.ErrEntityNotFound)
}

func TestGetWithErrorInDBReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()

	eavService.
		On("GetEntity", entityType.Name, uuid).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

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
		On("GetEntity", entityType.Name, uuid).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	entityReturned, err := controller.GetObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, entityReturned)
}

// ----------------------- GetAll -----------------------

func TestGetAllOfNotExistentTypeReturnsEmpty(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	eavService.
		On("GetEntities", "no-exists", map[string]string{}).
		Return([]*models.Entity{}, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/no-exists",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": "no-exists"})

	entities, err := controller.GetAll(response, request)
	assert.Nil(t, err)
	assert.Len(t, entities, 0)
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
		On("GetEntities", entityType.Name, map[string]string{}).
		Return([]*models.Entity{entity1, entity2}, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

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
		On("GetEntities", entityType.Name, map[string]string{"param1": "something"}).
		Return([]*models.Entity{entity1}, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
		"/objects/exists/",
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})
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

	uuid := uuid.New()

	eavService.
		On("DeleteEntity", entityType.Name, uuid).
		Return(errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"DELETE",
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	_, err := controller.DeleteObject(response, request)
	assert.ErrorContains(t, err, "db error")
}

func TestDeleteObjectReturnsNil(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

	entityType := &models.EntityType{
		Name: "entityType",
	}

	uuid := uuid.New()

	eavService.
		On("DeleteEntity", entityType.Name, uuid).
		Return(nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"DELETE",
		"/objects/exists/"+uuid.String(),
		strings.NewReader(""),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

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

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/objects/exists",
		strings.NewReader("bad json"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

	_, err := controller.CreateObject(response, request)
	assert.ErrorContains(t, err, "json decoding failed")
}

func TestCreteObjectWithErrorInDBReturnsError(t *testing.T) {
	eavService := mocksEAVService.NewEAVService(t)

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
		"POST",
		"/objects/exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

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
		On("CreateEntity", entityType.Name, map[string]any{"1": "1"}).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/objects/exists",
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name})

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

	eavService.
		On("UpdateEntity", entityType.Name, uuid, map[string]any{"1": "1"}).
		Return(nil, errors.New("db error"))

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"PUT",
		"/objects/exists/"+uuid.String(),
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

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
		On("UpdateEntity", entityType.Name, uuid, map[string]any{"1": "1"}).
		Return(entity, nil)

	controller := controllers.NewEAVController(
		logger,
		eavService,
	)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"PUT",
		"/objects/exists/"+uuid.String(),
		strings.NewReader("{\"1\": \"1\"}"),
	)
	request = mux.SetURLVars(request, map[string]string{"type": entityType.Name, "id": uuid.String()})

	responded, err := controller.ModifyObject(response, request)
	assert.Nil(t, err)
	assert.Equal(t, entity, responded)
}
