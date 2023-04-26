package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrEntityNotFound     = httperrors.NewErrorNotFound("entity", "please use a valid object id")
	ErrEntityTypeNotFound = httperrors.NewErrorNotFound("entity type", "please use a type that exists")
	ErrIDNotAnUUID        = httperrors.NewBadRequestError("id is not an uuid", "please use an uuid for the id value")
	ErrDBQueryFailed      = func(err error) httperrors.HTTPError {
		return httperrors.NewInternalServerError("db error", "database query failed", err)
	}
)

type EAVController interface {
	GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	UpdateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ EAVController = (*eavControllerImpl)(nil)

func NewEAVController(
	logger *zap.Logger,
	eavService services.EAVService,
) EAVController {
	return &eavControllerImpl{
		logger:     logger,
		eavService: eavService,
	}
}

// The concrete implementation of the EAVController
type eavControllerImpl struct {
	logger     *zap.Logger
	eavService services.EAVService
}

// The handler responsible of the retrieval of one objects
func (controller *eavControllerImpl) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, entityID, herr := controller.getEntityTypeNameAndEntityID(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.eavService.GetEntity(entityTypeName, entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	return entity, nil
}

// The handler responsible of the retrieval of multiple objects
func (controller *eavControllerImpl) GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	params, herr := controller.decodeJSON(r)
	if herr != nil {
		return nil, herr
	}

	entities, err := controller.eavService.GetEntities(entityTypeName, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	return entities, nil
}

// The handler responsible of the creation of a object
func (controller *eavControllerImpl) CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := controller.decodeJSON(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.eavService.CreateEntity(entityTypeName, attrs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	w.Header().Add("Location", buildLocationString(entity))
	w.WriteHeader(http.StatusCreated)

	return entity, nil
}

func buildLocationString(et *models.Entity) string {
	return fmt.Sprintf("/objects/%s/%s", et.EntityType.Name, et.ID.String())
}

// The handler responsible for the updates of one object
func (controller *eavControllerImpl) UpdateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, entityID, herr := controller.getEntityTypeNameAndEntityID(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := controller.decodeJSON(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.eavService.UpdateEntity(entityTypeName, entityID, attrs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	return entity, nil
}

// The handler responsible for the deletion of a object
func (controller *eavControllerImpl) DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, entityID, herr := controller.getEntityTypeNameAndEntityID(r)
	if herr != nil {
		return nil, herr
	}

	err := controller.eavService.DeleteEntity(entityTypeName, entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	return nil, nil
}

// Extract the "type" parameter from url
func (controller *eavControllerImpl) getEntityTypeNameFromRequest(r *http.Request) (string, httperrors.HTTPError) {
	entityTypeName, present := mux.Vars(r)["type"]
	if !present {
		return "", ErrEntityTypeNotFound
	}

	return entityTypeName, nil
}

// Extract the "id" parameter from url
func (controller *eavControllerImpl) getEntityIDFromRequest(r *http.Request) (uuid.UUID, httperrors.HTTPError) {
	id, present := mux.Vars(r)["id"]
	if !present {
		return uuid.Nil, ErrEntityNotFound
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, ErrIDNotAnUUID
	}

	return uid, nil
}

// Extract the "type" and "id" parameters from url
func (controller *eavControllerImpl) getEntityTypeNameAndEntityID(r *http.Request) (string, uuid.UUID, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return "", uuid.Nil, herr
	}

	entityID, herr := controller.getEntityIDFromRequest(r)
	if herr != nil {
		return "", uuid.Nil, herr
	}

	return entityTypeName, entityID, nil
}

// Decode json present in request body
func (controller *eavControllerImpl) decodeJSON(r *http.Request) (map[string]any, httperrors.HTTPError) {
	var attrs map[string]any
	err := json.NewDecoder(r.Body).Decode(&attrs)
	switch {
	case err == io.EOF:
		// empty body
		return map[string]any{}, nil
	case err != nil:
		return nil, httperrors.NewBadRequestError("json decoding failed", "please use a correct json payload")
	}

	return attrs, nil
}
