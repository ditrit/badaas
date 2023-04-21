package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	uuid "github.com/google/uuid"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/eavservice"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrEntityNotFound            = httperrors.NewErrorNotFound("entity", "please use a valid entity id")
	ErrEntityTypeNotFound        = httperrors.NewErrorNotFound("entity type", "please use a type that exists in the schema")
	ErrIDNotAnUUID               = httperrors.NewBadRequestError("id is not an uuid", "please use an uuid for the id value")
	ErrEntityTypeDontMatchEntity = httperrors.NewBadRequestError("types don't match", "the entity you are targeting don't match the entity type name provided")
	ErrDBQueryFailed             = func(err error) httperrors.HTTPError {
		return httperrors.NewInternalServerError("db error", "database query failed", err)
	}
)

type EAVController interface {
	GetAll(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	ModifyObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ EAVController = (*eavControllerImpl)(nil)

// The InformationController constructor
func NewEAVController(
	logger *zap.Logger,
	eavService eavservice.EAVService,
) EAVController {
	return &eavControllerImpl{
		logger:     logger,
		eavService: eavService,
	}
}

// The concrete implementation of the InformationController
type eavControllerImpl struct {
	logger     *zap.Logger
	eavService eavservice.EAVService
}

func (controller *eavControllerImpl) GetAll(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	ett, herr := controller.getEntityTypeFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	var qp = make(map[string]string)
	for k, v := range r.URL.Query() {
		qp[k] = v[0]
	}

	return controller.eavService.GetEntitiesWithParams(ett, qp), nil
}

// The handler responsible for the retrieval of une entity
func (controller *eavControllerImpl) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	return controller.getEntityFromRequest(r)
}

// The handler responsible for the deletion of entities and their associated value
func (controller *eavControllerImpl) DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entity, herr := controller.getEntityFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	err := controller.eavService.DeleteEntity(entity)
	if err != nil {
		return nil, ErrDBQueryFailed(err)
	}

	return nil, nil
}

// The handler responsible for the creation of entities
func (controller *eavControllerImpl) CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	ett, herr := controller.getEntityTypeFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := controller.decodeJSON(r)
	if herr != nil {
		return nil, herr
	}

	et, err := controller.eavService.CreateEntity(ett, attrs)
	if err != nil {
		return nil, ErrDBQueryFailed(err)
	}

	w.Header().Add("Location", buildLocationString(et))
	w.WriteHeader(http.StatusCreated)

	return et, nil
}

func buildLocationString(et *models.Entity) string {
	return fmt.Sprintf("/v1/objects/%s/%s", et.EntityType.Name, et.ID.String())
}

// The handler responsible for the updates of entities
func (controller *eavControllerImpl) ModifyObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entity, herr := controller.getEntityFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := controller.decodeJSON(r)
	if herr != nil {
		return nil, herr
	}

	err := controller.eavService.UpdateEntity(entity, attrs)
	if err != nil {
		return nil, ErrDBQueryFailed(err)
	}

	return entity, nil
}

// Extract the "type" parameter from url and get the entity type from the db
func (controller *eavControllerImpl) getEntityTypeFromRequest(r *http.Request) (*models.EntityType, httperrors.HTTPError) {
	entityTypeName, present := mux.Vars(r)["type"]
	if !present {
		return nil, ErrEntityTypeNotFound
	}

	ett, err := controller.eavService.GetEntityTypeByName(entityTypeName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}

	return ett, nil
}

// Gets the Entity object for the given ett and id or returns the appropriate http errors
func (controller *eavControllerImpl) getEntity(ett *models.EntityType, id uuid.UUID) (*models.Entity, httperrors.HTTPError) {
	entity, err := controller.eavService.GetEntity(ett, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		} else if errors.Is(err, eavservice.ErrIDDontMatchEntityType) {
			return nil, ErrEntityTypeDontMatchEntity
		} else {
			return nil, ErrDBQueryFailed(err)
		}
	}

	return entity, nil
}

// Gets the Entity object for the given ett and id or returns the appropriate http errors
func (controller *eavControllerImpl) getEntityFromRequest(r *http.Request) (*models.Entity, httperrors.HTTPError) {
	ett, herr := controller.getEntityTypeFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	id, herr := controller.getEntityIDFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	entity, herr := controller.getEntity(ett, id)
	if herr != nil {
		return nil, herr
	}

	return entity, nil
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

func (controller *eavControllerImpl) decodeJSON(r *http.Request) (map[string]any, httperrors.HTTPError) {
	var attrs map[string]any
	err := json.NewDecoder(r.Body).Decode(&attrs)
	if err != nil {
		return nil, httperrors.NewBadRequestError("json decoding failed", "please use a correct json payload")
	}

	return attrs, nil
}
