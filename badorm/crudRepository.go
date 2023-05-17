package badorm

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
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
	ErrObjectsNotRelated      = func(typeName, attributeName string) error {
		return fmt.Errorf("%[1]s has not attribute named %[2]s or %[2]sID", typeName, attributeName)
	}
	ErrModelNotRegistered = func(typeName, attributeName string) error {
		return fmt.Errorf("%[1]s has an attribute named %[2]s or %[2]sID but %[2]s is not registered as model (use AddModel)", typeName, attributeName)
	}
)

// Implementation of the Generic CRUD Repository
type CRUDRepositoryImpl[T any, ID BadaasID] struct {
	CRUDRepository[T, ID]
	logger *zap.Logger
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T any, ID BadaasID](
	logger *zap.Logger,
) CRUDRepository[T, ID] {
	return &CRUDRepositoryImpl[T, ID]{
		logger: logger,
	}
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
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
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
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
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
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
func (repository *CRUDRepositoryImpl[T, ID]) GetMultiple(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error) {
	thisEntityConditions, joinConditions := divideConditionsByEntity(conditions)

	query := tx.Where(getWhereParams(thisEntityConditions))

	initialTableName, err := getTableName(query, *new(T))
	if err != nil {
		return nil, err
	}

	for _, condition := range joinConditions {
		err := condition.ApplyTo(query, initialTableName)
		if err != nil {
			return nil, err
		}
	}

	// execute query
	var entities []*T
	err = query.Find(&entities).Error

	return entities, err
}

func getWhereParams[T any](conditions []WhereCondition[T]) map[string]any {
	// TODO verificar que no se repitan
	whereParams := map[string]any{}
	for _, condition := range conditions {
		whereParams[condition.Field] = condition.Value
	}

	return whereParams
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

// Returns an object of the type of the "entity" attribute called "relationName"
// and a boolean value indicating whether the id attribute that relates them
// in the database is in the "entity"'s table.
// Returns error if "entity" not a relation called "relationName".
// func getRelatedObject(entity any, relationName string) (any, bool, error) {
// 	entityType := getEntityType(entity)

// 	field, isPresent := entityType.FieldByName(relationName)
// 	if !isPresent {
// 		// some gorm relations dont have a direct relation in the model, only the id
// 		relatedObject, err := getRelatedObjectByID(entityType, relationName)
// 		if err != nil {
// 			return nil, false, err
// 		}

// 		return relatedObject, true, nil
// 	}

// 	_, isIDPresent := entityType.FieldByName(relationName + "ID")

// 	return createObject(field.Type), isIDPresent, nil
// }

// Returns an object of the type of the "entity" attribute called "relationName" + "ID"
// Returns error if "entity" not a relation called "relationName" + "ID"
// func getRelatedObjectByID(entityType reflect.Type, relationName string) (any, error) {
// 	_, isPresent := entityType.FieldByName(relationName + "ID")
// 	if !isPresent {
// 		return nil, ErrObjectsNotRelated(entityType.Name(), relationName)
// 	}

// 	// TODO foreignKey can be redefined (https://gorm.io/docs/has_one.html#Override-References)
// 	fieldType, isPresent := modelsMapping[relationName]
// 	if !isPresent {
// 		return nil, ErrModelNotRegistered(entityType.Name(), relationName)
// 	}

// 	return createObject(fieldType), nil
// }

// // Creates an object of type reflect.Type using reflection
// func createObject(entityType reflect.Type) any {
// 	return reflect.New(entityType).Elem().Interface()
// }
