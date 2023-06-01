package badorm

import (
	"gorm.io/gorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T any, ID BadaasID] interface {
	GetEntity(id ID) (*T, error)
	// GetEntities(conditions map[string]any) ([]*T, error)
	GetEntities(conditions ...Condition[T]) ([]*T, error)
}

// check interface compliance
var _ CRUDService[UUIDModel, UUID] = (*crudServiceImpl[UUIDModel, UUID])(nil)

// Implementation of the CRUD Service
type crudServiceImpl[T any, ID BadaasID] struct {
	CRUDService[T, ID]
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T any, ID BadaasID](
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		db:         db,
		repository: repository,
	}
}

// Get the object of type T that has the "id"
func (service *crudServiceImpl[T, ID]) GetEntity(id ID) (*T, error) {
	return service.repository.GetByID(service.db, id)
}

// Get entities of type T that match all "conditions"
func (service *crudServiceImpl[T, ID]) GetEntities(conditions ...Condition[T]) ([]*T, error) {
	return service.repository.GetMultiple(service.db, conditions...)
}
