package badorm

import (
	"gorm.io/gorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T Model, ID ModelID] interface {
	GetByID(id ID) (*T, error)
	Query(conditions ...Condition[T]) ([]*T, error)
}

// check interface compliance
var _ CRUDService[UUIDModel, UUID] = (*crudServiceImpl[UUIDModel, UUID])(nil)

// Implementation of the CRUD Service
type crudServiceImpl[T Model, ID ModelID] struct {
	CRUDService[T, ID]
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T Model, ID ModelID](
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		db:         db,
		repository: repository,
	}
}

// Get the object of type T that has the "id"
func (service *crudServiceImpl[T, ID]) GetByID(id ID) (*T, error) {
	return service.repository.GetByID(service.db, id)
}

// Get entities of type T that match all "conditions"
func (service *crudServiceImpl[T, ID]) Query(conditions ...Condition[T]) ([]*T, error) {
	return service.repository.Query(service.db, conditions...)
}
