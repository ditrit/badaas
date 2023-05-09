package controllers

import (
	"net/http"

	"github.com/ditrit/badaas/httperrors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// check interface compliance
var _ CRUDController = (*GeneralCRUDController)(nil)

func NewGeneralCRUDController(
	logger *zap.Logger,
	entityMapping map[string]CRUDController,
) *GeneralCRUDController {
	return &GeneralCRUDController{
		logger:        logger,
		entityMapping: entityMapping,
	}
}

// The concrete implementation of the EAVController
type GeneralCRUDController struct {
	logger        *zap.Logger
	entityMapping map[string]CRUDController
}

// The handler responsible of the retrieval of one objects
func (controller *GeneralCRUDController) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	crudController, herr := controller.getEntityCRUDController(r)
	if herr != nil {
		return nil, herr
	}

	return crudController.GetObject(w, r)
}

// The handler responsible of the retrieval of multiple objects
func (controller *GeneralCRUDController) GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	crudController, herr := controller.getEntityCRUDController(r)
	if herr != nil {
		return nil, herr
	}

	return crudController.GetObjects(w, r)
}

// The handler responsible of the creation of a object
func (controller *GeneralCRUDController) CreateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	crudController, herr := controller.getEntityCRUDController(r)
	if herr != nil {
		return nil, herr
	}

	return crudController.CreateObject(w, r)
}

// The handler responsible for the updates of one object
func (controller *GeneralCRUDController) UpdateObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	crudController, herr := controller.getEntityCRUDController(r)
	if herr != nil {
		return nil, herr
	}

	return crudController.UpdateObject(w, r)
}

// The handler responsible for the deletion of a object
func (controller *GeneralCRUDController) DeleteObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	crudController, herr := controller.getEntityCRUDController(r)
	if herr != nil {
		return nil, herr
	}

	return crudController.DeleteObject(w, r)
}

func (controller *GeneralCRUDController) getEntityCRUDController(r *http.Request) (CRUDController, httperrors.HTTPError) {
	entityTypeName, herr := controller.getEntityTypeNameFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	crudController, isPresent := controller.entityMapping[entityTypeName]
	if !isPresent {
		return nil, ErrEntityTypeNotFound
	}

	return crudController, nil
}

// Extract the "type" parameter from url
func (controller *GeneralCRUDController) getEntityTypeNameFromRequest(r *http.Request) (string, httperrors.HTTPError) {
	entityTypeName, present := mux.Vars(r)["type"]
	if !present {
		return "", ErrEntityTypeNotFound
	}

	return entityTypeName, nil
}
