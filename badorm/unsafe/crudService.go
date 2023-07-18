package unsafe

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T badorm.Model, ID badorm.ModelID] interface {
	GetEntities(conditions map[string]any) ([]*T, error)
}

// check interface compliance
var _ CRUDService[badorm.UUIDModel, badorm.UUID] = (*crudServiceImpl[badorm.UUIDModel, badorm.UUID])(nil)

// Implementation of the CRUD Service
type crudServiceImpl[T badorm.Model, ID badorm.ModelID] struct {
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDUnsafeService[T badorm.Model, ID badorm.ModelID](
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		db:         db,
		repository: repository,
	}
}

// Get entities of type T that match all "conditions"
// "params" is in {"attributeName": expectedValue} format
// in case of join "params" can have the format:
// {"relationAttributeName": {"attributeName": expectedValue}}
func (service *crudServiceImpl[T, ID]) GetEntities(conditions map[string]any) ([]*T, error) {
	return service.repository.GetMultiple(service.db, conditions)
}
