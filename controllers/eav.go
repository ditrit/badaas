package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	uuid "github.com/google/uuid"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrEntityNotFound     = httperrors.NewErrorNotFound("entity", "please use a valid entity id")
	ErrEntityTypeNotFound = httperrors.NewErrorNotFound("entity type", "please use a type that exists in the schema")
	ErrIDNotAnUUID        = httperrors.NewBadRequestError("id is not an uuid", "please use an uuid for the id value")
	ErrDBQueryFailed      = func(err error) httperrors.HTTPError {
		return httperrors.NewInternalServerError("db error", "database query failed", err)
	}
)

type EAVController interface {
	GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	GetAll(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	ModifyObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ EAVController = (*eavControllerImpl)(nil)

// The InformationController constructor
func NewEAVController(
	logger *zap.Logger,
	eavService services.EAVService,
) EAVController {
	return &eavControllerImpl{
		logger:     logger,
		eavService: eavService,
	}
}

// The concrete implementation of the InformationController
type eavControllerImpl struct {
	logger     *zap.Logger
	eavService services.EAVService
}

// The handler responsible for the retrieval of une entity
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

func (controller *eavControllerImpl) GetAll(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	var qp = make(map[string]string)
	for k, v := range r.URL.Query() {
		qp[k] = v[0]
	}

	entities, err := controller.eavService.GetEntities(entityTypeName, qp)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	return entities, nil
}

// The handler responsible for the creation of entities
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

// The handler responsible for the updates of entities
func (controller *eavControllerImpl) ModifyObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
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

// The handler responsible for the deletion of entities and their associated value
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

func (controller *eavControllerImpl) decodeJSON(r *http.Request) (map[string]any, httperrors.HTTPError) {
	var attrs map[string]any
	err := json.NewDecoder(r.Body).Decode(&attrs)
	if err != nil {
		return nil, httperrors.NewBadRequestError("json decoding failed", "please use a correct json payload")
	}

	return attrs, nil
}
