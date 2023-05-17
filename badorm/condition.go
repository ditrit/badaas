package badorm

import (
	"fmt"
	"reflect"

	"github.com/ettle/strcase"
	"gorm.io/gorm"
)

type Condition[T any] interface {
	ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error)

	// this method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	interfaceVerificationMethod(T)
}

type WhereCondition[T any] struct {
	Field string
	Value any
}

func (condition WhereCondition[T]) interfaceVerificationMethod(t T) {}

func (condition WhereCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	return query.Where(
		fmt.Sprintf("%s = ?", condition.Field),
		condition.Value,
	), nil
}

type JoinCondition[T1 any, T2 any] struct {
	Field      string
	Conditions []Condition[T2]
}

func (condition JoinCondition[T1, T2]) interfaceVerificationMethod(t T1) {}

func (condition JoinCondition[T1, T2]) ApplyTo(query *gorm.DB, previousTableName string) (*gorm.DB, error) {
	joinTableName, err := getTableName(query, *new(T2))
	if err != nil {
		return nil, err
	}

	tableWithSuffix := joinTableName + "_" + previousTableName

	var stringQuery string
	if isIDPresentInObject[T1](condition.Field) {
		stringQuery = fmt.Sprintf(
			`JOIN %[1]s %[2]s ON
				%[2]s.id = %[3]s.%[4]s_id
				AND %[2]s.deleted_at IS NULL
			`,
			joinTableName,
			tableWithSuffix,
			previousTableName,
			condition.Field,
		)
	} else {
		// TODO foreignKey can be redefined (https://gorm.io/docs/has_one.html#Override-References)
		previousAttribute := reflect.TypeOf(*new(T1)).Name()
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

	thisEntityConditions, joinConditions := divideConditionsByEntity(condition.Conditions)

	conditionsValues := []any{}
	for _, condition := range thisEntityConditions {
		stringQuery += fmt.Sprintf(
			`AND %[1]s.%[2]s = ?
			`,
			tableWithSuffix, condition.Field,
		)
		conditionsValues = append(conditionsValues, condition.Value)
	}

	query = query.Joins(stringQuery, conditionsValues...)

	for _, joinCondition := range joinConditions {
		query, err = joinCondition.ApplyTo(query, tableWithSuffix)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

func isIDPresentInObject[T any](relationName string) bool {
	entityType := getEntityType(*new(T))
	_, isIDPresent := entityType.FieldByName(
		strcase.ToPascal(relationName) + "ID",
	)
	return isIDPresent
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
func divideConditionsByEntity[T any](
	conditions []Condition[T],
) (thisEntityConditions []WhereCondition[T], joinConditions []Condition[T]) {
	for _, condition := range conditions {
		switch typedCondition := condition.(type) {
		case WhereCondition[T]:
			thisEntityConditions = append(thisEntityConditions, typedCondition)
		// case JoinCondition[T, any]:
		// joinConditions = append(joinConditions, typedCondition)
		default:
			joinConditions = append(joinConditions, typedCondition)
			// log.Println(reflect.TypeOf(typedCondition))
			// log.Println(reflect.TypeOf(any(condition)))
			// log.Println(condition.(JoinCondition[T, any]))
		}
	}

	return
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
