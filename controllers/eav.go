package controllers

import (
	"errors"
	"fmt"
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
)

// check interface compliance
var _ CRUDController = (*eavControllerImpl)(nil)

func NewEAVController(
	logger *zap.Logger,
	eavService services.EAVService,
) CRUDController {
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
		return nil, httperrors.NewDBError(err)
	}

	return entity, nil
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	return entities, nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	w.Header().Add("Location", buildLocationString(entity))
	w.WriteHeader(http.StatusCreated)

	return entity, nil
}

func buildLocationString(et *models.Entity) string {
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, httperrors.NewDBError(err)
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
		return nil, httperrors.NewDBError(err)
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

// Extract the "type" and "id" parameters from url
func (controller *eavControllerImpl) getEntityTypeNameAndEntityID(r *http.Request) (string, uuid.UUID, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return "", uuid.Nil, herr
	}

	entityID, herr := getEntityIDFromRequest(r)
	if herr != nil {
		return "", uuid.Nil, herr
	}

	return entityTypeName, entityID, nil
}
