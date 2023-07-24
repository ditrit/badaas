package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/unsafe"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
)

type CRUDController interface {
	// The handler responsible of the retrieval of one model
	GetModel(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)

	// The handler responsible of the retrieval of multiple models
	GetModels(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ CRUDController = (*crudControllerImpl[models.User])(nil)

type CRUDRoute struct {
	TypeName   string
	Controller CRUDController
}

func NewCRUDController[T badorm.Model](
	logger *zap.Logger,
	crudService badorm.CRUDService[T, badorm.UUID],
	unsafeCRUDService unsafe.CRUDService[T, badorm.UUID],
) CRUDRoute {
	fullTypeName := strings.ToLower(fmt.Sprintf("%T", *new(T)))
	// remove the package name of the type
	typeName := pie.Last(strings.Split(fullTypeName, "."))

	return CRUDRoute{
		TypeName: typeName,
		Controller: &crudControllerImpl[T]{
			logger:            logger,
			crudService:       crudService,
			unsafeCRUDService: unsafeCRUDService,
		},
	}
}

// The concrete implementation of the CRUDController
type crudControllerImpl[T badorm.Model] struct {
	logger            *zap.Logger
	crudService       badorm.CRUDService[T, badorm.UUID]
	unsafeCRUDService unsafe.CRUDService[T, badorm.UUID]
}

// The handler responsible of the retrieval of one model
func (controller *crudControllerImpl[T]) GetModel(_ http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityID, herr := getEntityIDFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.crudService.GetByID(entityID)

	return entity, mapServiceError(err)
}

// The handler responsible of the retrieval of multiple models
func (controller *crudControllerImpl[T]) GetModels(_ http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	params, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entities, err := controller.unsafeCRUDService.Query(params)

	return entities, mapServiceError(err)
}
