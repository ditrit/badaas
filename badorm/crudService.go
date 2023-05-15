package badorm

import (
	"github.com/ditrit/badaas/persistence/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CRUDService[T any, ID BadaasID] interface {
	GetEntity(id ID) (*T, error)
	GetEntities(conditions map[string]any) ([]*T, error)
}

// check interface compliance
var _ CRUDService[models.User, uuid.UUID] = (*crudServiceImpl[models.User, uuid.UUID])(nil)

// Implementation of the Generic CRUD Repository
type crudServiceImpl[T any, ID BadaasID] struct {
	CRUDService[T, ID]
	logger     *zap.Logger
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T any, ID BadaasID](
	logger *zap.Logger,
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

// Get the Entity of type with name "entityTypeName" that has the "id"
func (service *crudServiceImpl[T, ID]) GetEntity(id ID) (*T, error) {
	return service.repository.GetByID(service.db, id)
}

// Get entities of type with name "entityTypeName" that match all "conditions"
//
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
func (service *crudServiceImpl[T, ID]) GetEntities(conditions map[string]any) ([]*T, error) {
	return service.repository.GetMultiple(service.db, conditions)
}
