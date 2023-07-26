package badorm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm/logger"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T Model, ID ModelID] interface {
	// Get the model of type T that has the "id"
	GetByID(id ID) (*T, error)

	// Get only one model that match "conditions"
	// or returns error if 0 or more than 1 are found.
	QueryOne(conditions ...Condition[T]) (*T, error)

	// Get the list of models that match "conditions"
	Query(conditions ...Condition[T]) ([]*T, error)
}

// check interface compliance
var _ CRUDService[UUIDModel, UUID] = (*crudServiceImpl[UUIDModel, UUID])(nil)

// Implementation of the CRUD Service
type crudServiceImpl[T Model, ID ModelID] struct {
	CRUDService[T, ID]
	logger     logger.Interface
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T Model, ID ModelID](
	logger logger.Interface,
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		logger:     logger,
		db:         db,
		repository: repository,
	}
}

// Get the model of type T that has the "id"
func (service *crudServiceImpl[T, ID]) GetByID(id ID) (*T, error) {
	return Transaction[*T](
		service.logger,
		service.db,
		func(db *gorm.DB) (*T, error) {
			return service.repository.GetByID(db, id)
		},
	)
}

// Get only one model that match "conditions"
// or returns error if 0 or more than 1 are found.
func (service *crudServiceImpl[T, ID]) QueryOne(conditions ...Condition[T]) (*T, error) {
	return Transaction[*T](
		service.logger,
		service.db,
		func(db *gorm.DB) (*T, error) {
			return service.repository.QueryOne(db, conditions...)
		},
	)
}

// Get the list of models that match "conditions"
func (service *crudServiceImpl[T, ID]) Query(conditions ...Condition[T]) ([]*T, error) {
	return Transaction[[]*T](
		service.logger,
		service.db,
		func(db *gorm.DB) ([]*T, error) {
			return service.repository.Query(db, conditions...)
		},
	)
}
