package badorm

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Generic CRUD Repository
// T can be any model whose identifier attribute is of type ID
type CRUDUnsafeRepository[T any, ID BadaasID] interface {
	GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error)
}

var (
	ErrObjectsNotRelated = func(typeName, attributeName string) error {
		return fmt.Errorf("%[1]s has not attribute named %[2]s or %[2]sID", typeName, attributeName)
	}
	ErrModelNotRegistered = func(typeName, attributeName string) error {
		return fmt.Errorf(
			"%[1]s has an attribute named %[2]s or %[2]sID but %[2]s is not registered as model (use badorm.AddUnsafeModel or badorm.GetCRUDUnsafeServiceModule)",
			typeName, attributeName,
		)
	}
)

// Implementation of the Generic CRUD Repository
type CRUDUnsafeRepositoryImpl[T any, ID BadaasID] struct {
	CRUDUnsafeRepository[T, ID]
}

// Constructor of the Generic CRUD Repository
func NewCRUDUnsafeRepository[T any, ID BadaasID]() CRUDUnsafeRepository[T, ID] {
	return &CRUDUnsafeRepositoryImpl[T, ID]{}
}

// Get the list of objects that match "conditions" inside transaction "tx"
// "conditions" is in {"attributeName": expectedValue} format
// in case of join "conditions" can have the format:
// {"relationAttributeName": {"attributeName": expectedValue}}
func (repository *CRUDUnsafeRepositoryImpl[T, ID]) GetMultiple(tx *gorm.DB, conditions map[string]any) ([]*T, error) {
	thisEntityConditions, joinConditions, err := divideConditionsByEntity(conditions)
	if err != nil {
		return nil, err
	}

	// TODO on column definition, conditions need to have the column name in place of the attribute name
	query := tx.Where(thisEntityConditions)

	entity := new(T)
	// only entities that match the conditions
	for joinAttributeName, joinConditions := range joinConditions {
		var tableName string

		tableName, err = getTableName(tx, entity)
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

// Adds a join to the "query" by the "joinAttributeName"
// then, adds the verification that the joined values match "conditions"

// "conditions" is in {"attributeName": expectedValue} format
// "previousEntity" is a pointer to a object from where we navigate the relationship
// "previousTableName" is the name of the table where the previous object is saved and from we the join will we done
func (repository *CRUDUnsafeRepositoryImpl[T, ID]) addJoinToQuery(
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
		err = repository.addJoinToQuery(
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
		case string:
			uuid, err := ParseUUID(typedExpectedValue)
			if err == nil {
				thisEntityConditions[attributeName] = uuid
			} else {
				thisEntityConditions[attributeName] = expectedValue
			}
		case float64, bool, int:
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
