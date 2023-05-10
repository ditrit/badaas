package controllers

import (
	"errors"
	"net/http"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CRUDController interface {
	GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	UpdateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// TODO un monton de codigo duplicado

// check interface compliance
var _ CRUDController = (*crudControllerImpl[models.User])(nil)

func NewCRUDController[T models.Tabler](
	logger *zap.Logger,
	crudService services.CRUDService[T, uuid.UUID],
) CRUDController {
	return &crudControllerImpl[T]{
		logger:      logger,
		crudService: crudService,
	}
}

// The concrete implementation of the EAVController
type crudControllerImpl[T models.Tabler] struct {
	logger      *zap.Logger
	crudService services.CRUDService[T, uuid.UUID]
}

// The handler responsible of the retrieval of one objects
func (controller *crudControllerImpl[T]) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityID, herr := controller.getEntityIDFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.crudService.GetEntity(entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	return entity, nil
}

// The handler responsible of the retrieval of multiple objects
func (controller *crudControllerImpl[T]) GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	params, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entities, err := controller.crudService.GetEntities(params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	return entities, nil
}

// The handler responsible of the creation of a object
func (controller *crudControllerImpl[T]) CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	attrs, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.crudService.CreateEntity(attrs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	// TODO ver como hacer esto
	// w.Header().Add("Location", buildLocationString(entity))
	w.WriteHeader(http.StatusCreated)

	return entity, nil
}

// The handler responsible for the updates of one object
func (controller *crudControllerImpl[T]) UpdateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityID, herr := controller.getEntityIDFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	attrs, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.crudService.UpdateEntity(entityID, attrs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	return entity, nil
}

// The handler responsible for the deletion of a object
func (controller *crudControllerImpl[T]) DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityID, herr := controller.getEntityIDFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	err := controller.crudService.DeleteEntity(entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		return nil, httperrors.NewDBError(err)
	}

	return nil, nil
}

// TODO codigo duplicado

// Extract the "id" parameter from url
func (controller *crudControllerImpl[T]) getEntityIDFromRequest(r *http.Request) (uuid.UUID, httperrors.HTTPError) {
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
