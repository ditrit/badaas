package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/eavservice"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrEntityNotFound            = httperrors.NewErrorNotFound("entity", "please use a valid entity id")
	ErrEntityTypeNotFound        = httperrors.NewErrorNotFound("entity type", "please use a type that exists in the schema")
	ErrIdNotAnUUID               = httperrors.NewBadRequestError("id is not an uuid", "please use an uuid for the id value")
	ErrEntityTypeDontMatchEntity = httperrors.NewBadRequestError("types don't match", "the entity you are targeting don't match the entity type name provided")
	ErrDBQueryFailed             = func(err error) httperrors.HTTPError {
		return httperrors.NewInternalServerError("db error", "database query failed", err)
	}
)

// The EAV controller
type EAVController interface {
	// Return the badaas server information
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
	paginationConfiguration configuration.PaginationConfiguration,
	db *gorm.DB,
	eavService eavservice.EAVService,
) EAVController {
	return &eavControllerImpl{
		logger:                  logger,
		paginationConfiguration: paginationConfiguration,
		db:                      db,
		eavService:              eavService,
	}
}

// The concrete implementation of the InformationController
type eavControllerImpl struct {
	logger                  *zap.Logger
	paginationConfiguration configuration.PaginationConfiguration
	db                      *gorm.DB
	eavService              eavservice.EAVService
}

func (controller *eavControllerImpl) GetAll(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName := getEntityTypeNameFromRequest(r)

	ett, err := controller.eavService.GetEntityTypeByName(entityTypeName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}

	}
	queryparams := r.URL.Query()
	var qp = make(map[string]string)
	for k, v := range queryparams {
		qp[k] = v[0]
	}
	fmt.Println(qp)
	var collection = controller.eavService.GetEntitiesWithParams(ett, qp)

	return collection, nil
}

// The handler responsible for the retrieval of une entity
func (controller *eavControllerImpl) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {

	entityTypeName := getEntityTypeNameFromRequest(r)
	controller.logger.Sugar().Debugf("query for %s", entityTypeName)
	ett, err := controller.eavService.GetEntityTypeByName(entityTypeName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		}
		http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
		return nil, httperrors.NewInternalServerError("db error", "db query failed", err)
	}
	id, idErr := getEntityIDFromRequest(r)

	if idErr != nil {
		return nil, idErr
	}

	obj, err := controller.eavService.GetEntity(ett, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		} else if errors.Is(err, eavservice.ErrIdDontMatchEntityType) {
			http.Error(w, GetErrMsg(err.Error()), http.StatusBadRequest)
			return nil, ErrEntityTypeDontMatchEntity
		} else {
			return nil, ErrDBQueryFailed(err)
		}
	}
	return obj, nil
}

// Extract the "id" parameter from url
func getEntityIDFromRequest(r *http.Request) (uuid.UUID, httperrors.HTTPError) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return uuid.Nil, ErrEntityNotFound
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, ErrIdNotAnUUID
	}
	return uid, nil
}

// Extract the "type" parameter from url
func getEntityTypeNameFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	entityType := vars["type"]
	return entityType
}

// The handler responsible for the deletion of entities and their associated value
func (controller *eavControllerImpl) DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName := getEntityTypeNameFromRequest(r)
	ett, err := controller.eavService.GetEntityTypeByName(entityTypeName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
		return nil, httperrors.NewInternalServerError("db error", "search for entity type", err)
	}
	id, idErr := getEntityIDFromRequest(r)

	if idErr != nil {
		return nil, idErr
	}
	entity, err := controller.eavService.GetEntity(ett, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		} else if errors.Is(err, eavservice.ErrIdDontMatchEntityType) {
			return nil, ErrEntityTypeDontMatchEntity
		} else {
			return nil, ErrDBQueryFailed(err)
		}
	}
	err = controller.eavService.DeleteEntity(entity)
	if err != nil {
		return nil, httperrors.NewInternalServerError("deletion failed", "", err)
	}
	return nil, nil
}

// The handler responsible for the creation of entities
func (controller *eavControllerImpl) CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName := getEntityTypeNameFromRequest(r)
	ett, err := controller.eavService.GetEntityTypeByName(entityTypeName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, ErrDBQueryFailed(err)
	}
	var cr createReq
	err = json.NewDecoder(r.Body).Decode(&cr)
	r.Body.Close()
	if err != nil {
		return nil, httperrors.NewBadRequestError("json decoding failed", "please use a correct json payload for entity creation")
	}

	fmt.Println(cr)
	et, err := controller.eavService.CreateEntity(ett, cr.Attrs)
	if err != nil {
		return nil, httperrors.NewInternalServerError("creation failed", "", err)
	}
	w.Header().Add("Location", buildLocationString(et))
	return et, nil
}

func buildLocationString(et *models.Entity) string {
	return fmt.Sprintf("/v1/objects/%s/%d", et.EntityType.Name, et.ID)
}

type createReq struct {
	Attrs map[string]interface{}
}

// The handler responsible for the updates of entities
func (controller *eavControllerImpl) ModifyObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityTypeName := getEntityTypeNameFromRequest(r)
	ett, err := controller.eavService.GetEntityTypeByName(entityTypeName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityTypeNotFound
		}
		return nil, ErrDBQueryFailed(err)

	}

	id, err := getEntityIDFromRequest(r)
	if err != nil {
		return nil, ErrIdNotAnUUID
	}
	entity, err := controller.eavService.GetEntity(ett, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEntityNotFound
		} else if errors.Is(err, eavservice.ErrIdDontMatchEntityType) {
			return nil, ErrEntityTypeDontMatchEntity
		} else {
			return nil, ErrDBQueryFailed(err)
		}
	}

	var mr modifyReq
	err = json.NewDecoder(r.Body).Decode(&mr)
	if err != nil {
		return nil, ErrDBQueryFailed(err)
	}
	fmt.Println(mr.Attrs)
	controller.eavService.UpdateEntity(entity, mr.Attrs)
	return entity, nil
}

// return json formatted string to be consumed by frontend or client
func GetErrMsg(msg string) string {
	return fmt.Sprintf(
		`{"error": %q}`,
		msg,
	)
}

type modifyReq struct {
	Attrs map[string]interface{}
}
