package badorm

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T any, ID BadaasID] interface {
	GetEntity(id ID) (*T, error)
	// GetEntities(conditions map[string]any) ([]*T, error)
	GetEntities(conditions ...Condition[T]) ([]*T, error)
}

// check interface compliance
var _ CRUDService[UUIDModel, uuid.UUID] = (*crudServiceImpl[UUIDModel, uuid.UUID])(nil)

// Implementation of the CRUD Service
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

// Get the object of type T that has the "id"
func (service *crudServiceImpl[T, ID]) GetEntity(id ID) (*T, error) {
	return service.repository.GetByID(service.db, id)
}

// Get entities of type T that match all "conditions"
//
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
//
// TODO
// func (service *crudServiceImpl[T, ID]) GetEntities(conditions map[string]any) ([]*T, error) {
func (service *crudServiceImpl[T, ID]) GetEntities(conditions ...Condition[T]) ([]*T, error) {
	return service.repository.GetMultiple(service.db, conditions...)
}
