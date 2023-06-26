package badorm

import (
	"gorm.io/gorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDUnsafeService[T Model, ID ModelID] interface {
	GetEntities(conditions map[string]any) ([]*T, error)
}

// check interface compliance
var _ CRUDUnsafeService[UUIDModel, UUID] = (*crudUnsafeServiceImpl[UUIDModel, UUID])(nil)

// Implementation of the CRUD Service
type crudUnsafeServiceImpl[T Model, ID ModelID] struct {
	CRUDService[T, ID]
	db         *gorm.DB
	repository CRUDUnsafeRepository[T, ID]
}

func NewCRUDUnsafeService[T Model, ID ModelID](
	db *gorm.DB,
	repository CRUDUnsafeRepository[T, ID],
) CRUDUnsafeService[T, ID] {
	return &crudUnsafeServiceImpl[T, ID]{
		db:         db,
		repository: repository,
	}
}

// Get entities of type T that match all "conditions"
// "params" is in {"attributeName": expectedValue} format
// in case of join "params" can have the format:
// {"relationAttributeName": {"attributeName": expectedValue}}
func (service *crudUnsafeServiceImpl[T, ID]) GetEntities(conditions map[string]any) ([]*T, error) {
	return service.repository.GetMultiple(service.db, conditions)
}
