package badorm

import (
	"errors"

	"gorm.io/gorm"
)

// Generic CRUD Repository
// T can be any model whose identifier attribute is of type ID
type CRUDRepository[T Model, ID ModelID] interface {
	// Create model "model" inside transaction "tx"
	Create(tx *gorm.DB, entity *T) error

	// ----- read -----

	// Get a model by its ID
	GetByID(tx *gorm.DB, id ID) (*T, error)

	// Get only one model that match "conditions" inside transaction "tx"
	// or returns error if 0 or more than 1 are found.
	QueryOne(tx *gorm.DB, conditions ...Condition[T]) (*T, error)

	// Get the list of models that match "conditions" inside transaction "tx"
	Query(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error)

	// Save model "model" inside transaction "tx"
	Save(tx *gorm.DB, entity *T) error

	// Delete model "model" inside transaction "tx"
	Delete(tx *gorm.DB, entity *T) error
}

var (
	ErrMoreThanOneObjectFound = errors.New("found more that one object that meet the requested conditions")
	ErrObjectNotFound         = errors.New("no object exists that meets the requested conditions")
)

// Implementation of the Generic CRUD Repository
type crudRepositoryImpl[T Model, ID ModelID] struct {
	CRUDRepository[T, ID]
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T Model, ID ModelID]() CRUDRepository[T, ID] {
	return &crudRepositoryImpl[T, ID]{}
}

// Create model "model" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Create(tx *gorm.DB, model *T) error {
	return tx.Create(model).Error
}

// Delete model "model" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Delete(tx *gorm.DB, model *T) error {
	return tx.Delete(model).Error
}

// Save model "model" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Save(tx *gorm.DB, model *T) error {
	return tx.Save(model).Error
}

// Get a model by its ID
func (repository *crudRepositoryImpl[T, ID]) GetByID(tx *gorm.DB, id ID) (*T, error) {
	var model T

	err := tx.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Get only one model that match "conditions" inside transaction "tx"
// or returns error if 0 or more than 1 are found.
func (repository *crudRepositoryImpl[T, ID]) QueryOne(tx *gorm.DB, conditions ...Condition[T]) (*T, error) {
	entities, err := repository.Query(tx, conditions...)
	if err != nil {
		return nil, err
	}

	switch {
	case len(entities) == 1:
		return entities[0], nil
	case len(entities) == 0:
		return nil, ErrObjectNotFound
	default:
		return nil, ErrMoreThanOneObjectFound
	}
}

// Get the list of models that match "conditions" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Query(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error) {
	query, err := NewQuery(tx, conditions)
	if err != nil {
		return nil, err
	}

	// execute query
	var entities []*T
	err = query.Find(&entities)

	return entities, err
}
