package badorm

import (
	"gorm.io/gorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDUnsafeService[T any, ID BadaasID] interface {
	GetEntities(conditions map[string]any) ([]*T, error)
}

// check interface compliance
var _ CRUDUnsafeService[UUIDModel, UUID] = (*crudUnsafeServiceImpl[UUIDModel, UUID])(nil)

// Implementation of the CRUD Service
type crudUnsafeServiceImpl[T any, ID BadaasID] struct {
	CRUDService[T, ID]
	db         *gorm.DB
	repository CRUDUnsafeRepository[T, ID]
}

func NewCRUDUnsafeService[T any, ID BadaasID](
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
