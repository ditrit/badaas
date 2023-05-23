package badorm

import (
	"fmt"
	"reflect"

	"github.com/ettle/strcase"
	"gorm.io/gorm"
)

type Condition[T any] interface {
	// Applies the condition to the "query"
	// using the "tableName" as name for the table holding
	// the data for object of type T
	ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error)

	// This method is necessary to get the compiler to verify
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

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition WhereCondition[T]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	sql, values := condition.GetSQL(tableName)
	return query.Where(
		sql,
		values...,
	), nil
}

func (condition WhereCondition[T]) GetSQL(tableName string) (string, []any) {
	val := condition.Value
	// avoid nil is not nil behavior of go
	if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
		return fmt.Sprintf(
			"%s.%s IS NULL",
			tableName,
			condition.Field,
		), []any{}
	}

	return fmt.Sprintf(
		"%s.%s = ?",
		tableName,
		condition.Field,
	), []any{val}
}

type JoinCondition[T1 any, T2 any] struct {
	Field      string
	Conditions []Condition[T2]
}

func (condition JoinCondition[T1, T2]) interfaceVerificationMethod(t T1) {}

// Applies a join between the tables of T1 and T2
// previousTableName is the name of the table of T1
// It also applies the nested conditions
func (condition JoinCondition[T1, T2]) ApplyTo(query *gorm.DB, previousTableName string) (*gorm.DB, error) {
	// get the name of the table for T2
	toBeJoinedTableName, err := getTableName(query, *new(T2))
	if err != nil {
		return nil, err
	}

	// add a suffix to avoid tables with the same name when joining
	// the same table more than once
	nextTableName := toBeJoinedTableName + "_" + previousTableName

	// get the sql to do the join with T2
	joinQuery := condition.getSQLJoin(toBeJoinedTableName, nextTableName, previousTableName)

	whereConditions, joinConditions := divideConditionsByType(condition.Conditions)

	// apply WhereConditions to join in "on" clause
	conditionsValues := []any{}
	for _, condition := range whereConditions {
		sql, values := condition.GetSQL(nextTableName)
		joinQuery += "AND " + sql
		conditionsValues = append(conditionsValues, values...)
	}

	// add the join to the query
	query = query.Joins(joinQuery, conditionsValues...)

	// apply nested joins
	for _, joinCondition := range joinConditions {
		query, err = joinCondition.ApplyTo(query, nextTableName)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

// Returns the SQL string to do a join between T1 and T2
// taking into account that the ID attribute necessary to do it
// can be either in T1's or T2's table.
func (condition JoinCondition[T1, T2]) getSQLJoin(toBeJoinedTableName, nextTableName, previousTableName string) string {
	if isIDPresentInObject[T1](condition.Field) {
		// T1 has the id attribute
		return fmt.Sprintf(
			`JOIN %[1]s %[2]s ON
				%[2]s.id = %[3]s.%[4]s_id
				AND %[2]s.deleted_at IS NULL
			`,
			toBeJoinedTableName,
			nextTableName,
			previousTableName,
			condition.Field,
		)
	}
	// T2 has the id attribute
	// TODO foreignKey can be redefined (https://gorm.io/docs/has_one.html#Override-References)
	previousAttribute := reflect.TypeOf(*new(T1)).Name()
	return fmt.Sprintf(
		`JOIN %[1]s %[2]s ON
			%[2]s.%[4]s_id = %[3]s.id
			AND %[2]s.deleted_at IS NULL
		`,
		toBeJoinedTableName,
		nextTableName,
		previousTableName,
		previousAttribute,
	)
}

// Returns true if object T has an attribute called "relationName"ID
func isIDPresentInObject[T any](relationName string) bool {
	entityType := getEntityType(*new(T))
	_, isIDPresent := entityType.FieldByName(
		strcase.ToPascal(relationName) + "ID",
	)
	return isIDPresent
}

// Divides a list of conditions by its type: WhereConditions and JoinConditions
func divideConditionsByType[T any](
	conditions []Condition[T],
) (thisEntityConditions []WhereCondition[T], joinConditions []Condition[T]) {
	for _, condition := range conditions {
		switch typedCondition := condition.(type) {
		case WhereCondition[T]:
			thisEntityConditions = append(thisEntityConditions, typedCondition)
		default:
			joinConditions = append(joinConditions, typedCondition)
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
