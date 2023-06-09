package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elliotchance/pie/v2"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
)

type CRUDController interface {
	GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
	GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ CRUDController = (*crudControllerImpl[models.User])(nil)

type CRUDRoute struct {
	TypeName   string
	Controller CRUDController
}

func NewCRUDController[T any](
	logger *zap.Logger,
	crudService badorm.CRUDService[T, badorm.UUID],
	crudUnsafeService badorm.CRUDUnsafeService[T, badorm.UUID],
) CRUDRoute {
	fullTypeName := strings.ToLower(fmt.Sprintf("%T", *new(T)))
	// remove the package name of the type
	typeName := pie.Last(strings.Split(fullTypeName, "."))

	return CRUDRoute{
		TypeName: typeName,
		Controller: &crudControllerImpl[T]{
			logger:            logger,
			crudService:       crudService,
			crudUnsafeService: crudUnsafeService,
		},
	}
}

// The concrete implementation of the CRUDController
type crudControllerImpl[T any] struct {
	logger            *zap.Logger
	crudService       badorm.CRUDService[T, badorm.UUID]
	crudUnsafeService badorm.CRUDUnsafeService[T, badorm.UUID]
}

// The handler responsible of the retrieval of one object
func (controller *crudControllerImpl[T]) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityID, herr := getEntityIDFromRequest(r)
	if herr != nil {
		return nil, herr
	}

	entity, err := controller.crudService.GetEntity(entityID)

	return entity, mapServiceError(err)
}

// The handler responsible of the retrieval of multiple objects
func (controller *crudControllerImpl[T]) GetObjects(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	params, herr := decodeJSONOptional(r)
	if herr != nil {
		return nil, herr
	}

	entities, err := controller.crudUnsafeService.GetEntities(params)

	return entities, mapServiceError(err)
}
