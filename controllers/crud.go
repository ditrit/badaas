package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/elliotchance/pie/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	crudService badorm.CRUDService[T, uuid.UUID],
) CRUDRoute {
	fullTypeName := strings.ToLower(fmt.Sprintf("%T", *new(T)))
	// remove the package name of the type
	typeName := pie.Last(strings.Split(fullTypeName, "."))

	return CRUDRoute{
		TypeName: typeName,
		Controller: &crudControllerImpl[T]{
			logger:      logger,
			crudService: crudService,
		},
	}

}

// The concrete implementation of the EAVController
type crudControllerImpl[T any] struct {
	logger      *zap.Logger
	crudService badorm.CRUDService[T, uuid.UUID]
}

// The handler responsible of the retrieval of one objects
func (controller *crudControllerImpl[T]) GetObject(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	entityID, herr := getEntityIDFromRequest(r)
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
