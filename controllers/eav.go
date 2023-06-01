package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	return entity, mapServiceError(err)
}

// The handler responsible of the retrieval of multiple objects
func (controller *eavControllerImpl) GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	params, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entities, err := controller.eavService.GetEntities(entityTypeName, params)
	return entities, mapEAVServiceError(err)
}

// The handler responsible of the creation of a object
func (controller *eavControllerImpl) CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.eavService.CreateEntity(entityTypeName, attrs)
	if err != nil {
		return nil, mapEAVServiceError(err)
	}

	w.Header().Add("Location", buildEAVLocationString(entity))
	w.WriteHeader(http.StatusCreated)

	return entity, nil
}

func buildEAVLocationString(et *models.Entity) string {
	return fmt.Sprintf("eav/objects/%s/%s", et.EntityType.Name, et.ID.String())
}

// The handler responsible for the updates of one object
func (controller *eavControllerImpl) UpdateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, entityID, herr := controller.getEntityTypeNameAndEntityID(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.eavService.UpdateEntity(entityTypeName, entityID, attrs)
	return entity, mapServiceError(err)
}

// The handler responsible for the deletion of a object
func (controller *eavControllerImpl) DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName, entityID, herr := controller.getEntityTypeNameAndEntityID(r)
	if herr != nil {
		return nil, herr
	}

	err := controller.eavService.DeleteEntity(entityTypeName, entityID)
	return nil, mapServiceError(err)
}

func mapEAVServiceError(err error) httperrors.HTTPError {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEntityTypeNotFound
		}

		return mapServiceError(err)
	}

	return nil
}

// Extract the "type" parameter from url
func (controller *eavControllerImpl) getEntityTypeNameFromRequest(r *http.Request) (string, httperrors.HTTPError) {
	entityTypeName, present := mux.Vars(r)["type"]
	if !present {
		return "", ErrEntityTypeNotFound
	}

	return entityTypeName, nil
}

// Extract the "type" and "id" parameters from url
func (controller *eavControllerImpl) getEntityTypeNameAndEntityID(r *http.Request) (string, badorm.UUID, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return "", badorm.UUID(uuid.Nil), herr
	}

	entityID, herr := getEntityIDFromRequest(r)
	if herr != nil {
		return "", badorm.UUID(uuid.Nil), herr
	}

	return entityTypeName, entityID, nil
}
