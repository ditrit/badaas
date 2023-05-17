package badorm

import (
	"errors"
	"fmt"
	"reflect"
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
	Get(tx *gorm.DB, conditions map[string]any) (*T, error)
	GetOptional(tx *gorm.DB, conditions map[string]any) (*T, error)
	GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error)
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

// Get an object that matches "conditions" inside transaction "tx"
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
func (repository *CRUDRepositoryImpl[T, ID]) Get(tx *gorm.DB, conditions map[string]any) (*T, error) {
	entity, err := repository.GetOptional(tx, conditions)
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
func (repository *CRUDRepositoryImpl[T, ID]) GetOptional(tx *gorm.DB, conditions map[string]any) (*T, error) {
	entities, err := repository.GetMultiple(tx, conditions)
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
func (repository *CRUDRepositoryImpl[T, ID]) GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error) {
	thisEntityConditions, joinConditions, err := divideConditionsByEntity(conditions)
	if err != nil {
		return nil, err
	}

	query := tx.Where(thisEntityConditions)

	entity := new(T)
	// only entities that match the conditions
	for joinAttributeName, joinConditions := range joinConditions {
		tableName, err := getTableName(tx, entity)
		if err != nil {
			return nil, err
		}

		err = repository.addJoinToQuery(
			query,
			entity,
			tableName,
			joinAttributeName,
			joinConditions,
		)
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
	return repository.GetMultiple(tx, map[string]any{})
}

// Adds a join to the "query" by the "joinAttributeName"
// then, adds the verification that the joined values match "conditions"

// "conditions" is in {"attributeName": expectedValue} format
// "previousEntity" is a pointer to a object from where we navigate the relationship
// "previousTableName" is the name of the table where the previous object is saved and from we the join will we done
func (repository *CRUDRepositoryImpl[T, ID]) addJoinToQuery(
	query *gorm.DB, previousEntity any,
	previousTableName, joinAttributeName string,
	conditions map[string]any,
) error {
	thisEntityConditions, joinConditions, err := divideConditionsByEntity(conditions)
	if err != nil {
		return err
	}

	relatedObject, relationIDIsInPreviousTable, err := getRelatedObject(
		previousEntity,
		joinAttributeName,
	)
	if err != nil {
		return err
	}

	joinTableName, err := getTableName(query, relatedObject)
	if err != nil {
		return err
	}

	tableWithSuffix := joinTableName + "_" + previousTableName

	var stringQuery string
	if relationIDIsInPreviousTable {
		stringQuery = fmt.Sprintf(
			`JOIN %[1]s %[2]s ON
				%[2]s.id = %[3]s.%[4]s_id
				AND %[2]s.deleted_at IS NULL
			`,
			joinTableName,
			tableWithSuffix,
			previousTableName,
			joinAttributeName,
		)
	} else {
		// TODO foreignKey can be redefined (https://gorm.io/docs/has_one.html#Override-References)
		previousAttribute := reflect.TypeOf(previousEntity).Elem().Name()
		stringQuery = fmt.Sprintf(
			`JOIN %[1]s %[2]s ON
				%[2]s.%[4]s_id = %[3]s.id
				AND %[2]s.deleted_at IS NULL
			`,
			joinTableName,
			tableWithSuffix,
			previousTableName,
			previousAttribute,
		)
	}

	conditionsValues := []any{}
	for attributeName, conditionValue := range thisEntityConditions {
		stringQuery += fmt.Sprintf(
			`AND %[1]s.%[2]s = ?
			`,
			tableWithSuffix, attributeName,
		)
		conditionsValues = append(conditionsValues, conditionValue)
	}

	query.Joins(stringQuery, conditionsValues...)

	for joinAttributeName, joinConditions := range joinConditions {
		err := repository.addJoinToQuery(
			query,
			relatedObject,
			tableWithSuffix,
			joinAttributeName,
			joinConditions,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Given a map of "conditions" that is in {"attributeName": expectedValue} format
// and in case of join "conditions" can have the format:
//
//	{"relationAttributeName": {"attributeName": expectedValue}}
//
// it divides the map in two:
// the conditions that will be applied to the current entity ({"attributeName": expectedValue} format)
// the conditions that will generate a join with another entity ({"relationAttributeName": {"attributeName": expectedValue}} format)
//
// Returns error if any expectedValue is not of a supported type
func divideConditionsByEntity(
	conditions map[string]any,
) (map[string]any, map[string]map[string]any, error) {
	thisEntityConditions := map[string]any{}
	joinConditions := map[string]map[string]any{}

	for attributeName, expectedValue := range conditions {
		switch typedExpectedValue := expectedValue.(type) {
		case float64, bool, string, int:
			thisEntityConditions[attributeName] = expectedValue
		case map[string]any:
			joinConditions[attributeName] = typedExpectedValue
		default:
			return nil, nil, fmt.Errorf("unsupported type")
		}
	}

	return thisEntityConditions, joinConditions, nil
}

// Returns an object of the type of the "entity" attribute called "relationName"
// and a boolean value indicating whether the id attribute that relates them
// in the database is in the "entity"'s table.
// Returns error if "entity" not a relation called "relationName".
func getRelatedObject(entity any, relationName string) (any, bool, error) {
	entityType := getEntityType(entity)

	field, isPresent := entityType.FieldByName(relationName)
	if !isPresent {
		// some gorm relations dont have a direct relation in the model, only the id
		relatedObject, err := getRelatedObjectByID(entityType, relationName)
		if err != nil {
			return nil, false, err
		}

		return relatedObject, true, nil
	}

	_, isIDPresent := entityType.FieldByName(relationName + "ID")

	return createObject(field.Type), isIDPresent, nil
}

// Get the reflect.Type of any entity or pointer to entity
func getEntityType(entity any) reflect.Type {
	entityType := reflect.TypeOf(entity)

	// entityType will be a pointer if the relation can be nullable
	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
	}

	return entityType
}

// Returns an object of the type of the "entity" attribute called "relationName" + "ID"
// Returns error if "entity" not a relation called "relationName" + "ID"
func getRelatedObjectByID(entityType reflect.Type, relationName string) (any, error) {
	_, isPresent := entityType.FieldByName(relationName + "ID")
	if !isPresent {
		return nil, ErrObjectsNotRelated(entityType.Name(), relationName)
	}

	// TODO foreignKey can be redefined (https://gorm.io/docs/has_one.html#Override-References)
	fieldType, isPresent := modelsMapping[relationName]
	if !isPresent {
		return nil, ErrModelNotRegistered(entityType.Name(), relationName)
	}

	return createObject(fieldType), nil
}

// Creates an object of type reflect.Type using reflection
func createObject(entityType reflect.Type) any {
	return reflect.New(entityType).Elem().Interface()
}
