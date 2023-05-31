package badorm

import (
	"errors"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Generic CRUD Repository
// T can be any model whose identifier attribute is of type ID
type CRUDRepository[T any, ID BadaasID] interface {
	// create
	Create(tx *gorm.DB, entity *T) error
	// read
	GetByID(tx *gorm.DB, id ID) (*T, error)
	GetOptionalByID(tx *gorm.DB, id ID) (*T, error)
	Get(tx *gorm.DB, conditions ...Condition[T]) (*T, error)
	GetOptional(tx *gorm.DB, conditions ...Condition[T]) (*T, error)
	GetMultiple(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error)
	GetAll(tx *gorm.DB) ([]*T, error)
	// update
	Save(tx *gorm.DB, entity *T) error
	// delete
	Delete(tx *gorm.DB, entity *T) error
}

var (
	ErrMoreThanOneObjectFound = errors.New("found more that one object that meet the requested conditions")
	ErrObjectNotFound         = errors.New("no object exists that meets the requested conditions")
)

// Implementation of the Generic CRUD Repository
type CRUDRepositoryImpl[T any, ID BadaasID] struct {
	CRUDRepository[T, ID]
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T any, ID BadaasID]() CRUDRepository[T, ID] {
	return &CRUDRepositoryImpl[T, ID]{}
}

// Create object "entity" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Create(tx *gorm.DB, entity *T) error {
	return tx.Create(entity).Error
}

// Delete object "entity" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Delete(tx *gorm.DB, entity *T) error {
	return tx.Delete(entity).Error
}

// Save object "entity" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Save(tx *gorm.DB, entity *T) error {
	return tx.Save(entity).Error
}

// Get an object by "id" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) GetByID(tx *gorm.DB, id ID) (*T, error) {
	var entity T
	err := tx.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Get an object by "id" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) GetOptionalByID(tx *gorm.DB, id ID) (*T, error) {
	entity, err := repository.GetByID(tx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return entity, nil
}

// Get an object that matches "conditions" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Get(tx *gorm.DB, conditions ...Condition[T]) (*T, error) {
	entity, err := repository.GetOptional(tx, conditions...)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, ErrObjectNotFound
	}

	return entity, nil
}

// Get an object or nil that matches "conditions" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) GetOptional(tx *gorm.DB, conditions ...Condition[T]) (*T, error) {
	entities, err := repository.GetMultiple(tx, conditions...)
	if err != nil {
		return nil, err
	}

	if len(entities) > 1 {
		return nil, ErrMoreThanOneObjectFound
	} else if len(entities) == 1 {
		return entities[0], nil
	}

	return nil, nil
}

// Get the list of objects that match "conditions" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) GetMultiple(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error) {
	initialTableName, err := getTableName(tx, *new(T))
	if err != nil {
		return nil, err
	}

	query := tx
	for _, condition := range conditions {
		query, err = condition.ApplyTo(query, initialTableName)
		if err != nil {
			return nil, err
		}
	}

	// execute query
	var entities []*T
	err = query.Find(&entities).Error

	return entities, err
}

// Get the name of the table in "db" in which the data for "entity" is saved
// returns error is table name can not be found by gorm,
// probably because the type of "entity" is not registered using AddModel
func getTableName(db *gorm.DB, entity any) (string, error) {
	schemaName, err := schema.Parse(entity, &sync.Map{}, db.NamingStrategy)
	if err != nil {
		return "", err
	}

	return schemaName.Table, nil
}

// Get the list of objects of type T
func (repository *CRUDRepositoryImpl[T, ID]) GetAll(tx *gorm.DB) ([]*T, error) {
	return repository.GetMultiple(tx)
}
