package orm

import (
	"fmt"

	"gorm.io/gorm"
)

const DeletedAtField = "DeletedAt"

type ConditionType string

const (
	WhereType ConditionType = "Where"
	JoinType  ConditionType = "Join"
)

type Condition[T any] interface {
	// Applies the condition to the "query"
	// using the "tableName" as name for the table holding
	// the data for object of type T
	ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error)

	Type() ConditionType

	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T],
	// since if no method receives by parameter a type T,
	// any other Condition[T2] would also be considered a Condition[T].
	interfaceVerificationMethod(T)
}

type WhereCondition[T any] interface {
	GetSQL(query *gorm.DB, tableName string) (string, []any)

	getField() string
}

type FieldCondition[TObject any, TAtribute any] struct {
	Field        string
	Column       string
	ColumnPrefix string
	Expressions  []Expression[TAtribute]
}

//nolint:unused // see inside
func (condition FieldCondition[TObject, TAtribute]) interfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition FieldCondition[TObject, TAtribute]) ApplyTo(query *gorm.DB, tableName string) (*gorm.DB, error) {
	sql, values := condition.GetSQL(query, tableName)

	if condition.Field == DeletedAtField {
		query = query.Unscoped()
	}

	return query.Where(
		sql,
		values...,
	), nil
}

func (condition FieldCondition[TObject, TAtribute]) Type() ConditionType {
	return WhereType
}

func (condition FieldCondition[TObject, TAtribute]) getField() string {
	return condition.Field
}

func (condition FieldCondition[TObject, TAtribute]) GetSQL(query *gorm.DB, tableName string) (string, []any) {
	columnName := condition.Column
	if columnName == "" {
		columnName = query.NamingStrategy.ColumnName(tableName, condition.Field)
	}
	columnName = condition.ColumnPrefix + columnName

	conditionString := ""
	values := []any{}

	for index, expression := range condition.Expressions {
		// TODO que se pueda hacer la connection distinta aca
		// TODO strings.Join(exprs, " AND ")?
		if index != 0 {
			conditionString += " AND "
		}

		expressionSQL, expressionValues := expression.ToSQL(
			fmt.Sprintf(
				"%s.%s",
				tableName,
				columnName,
			),
		)
		conditionString += expressionSQL

		values = append(values, expressionValues...)
	}

	return conditionString, values
}

type JoinCondition[T1 any, T2 any] struct {
	T1Field    string
	T2Field    string
	Conditions []Condition[T2]
}

func (condition JoinCondition[T1, T2]) interfaceVerificationMethod(t T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

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
	joinQuery := condition.getSQLJoin(query, toBeJoinedTableName, nextTableName, previousTableName)

	whereConditions, joinConditions := divideConditionsByType(condition.Conditions)

	// apply WhereConditions to join in "on" clause
	conditionsValues := []any{}
	isDeletedAtConditionPresent := false
	for _, condition := range whereConditions {
		if condition.getField() == DeletedAtField {
			isDeletedAtConditionPresent = true
		}
		sql, values := condition.GetSQL(query, nextTableName)
		joinQuery += " AND " + sql
		conditionsValues = append(conditionsValues, values...)
	}

	if !isDeletedAtConditionPresent {
		joinQuery += fmt.Sprintf(
			" AND %s.deleted_at IS NULL",
			nextTableName,
		)
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

func (condition JoinCondition[T1, T2]) Type() ConditionType {
	return JoinType
}

// Returns the SQL string to do a join between T1 and T2
// taking into account that the ID attribute necessary to do it
// can be either in T1's or T2's table.
func (condition JoinCondition[T1, T2]) getSQLJoin(query *gorm.DB, toBeJoinedTableName, nextTableName, previousTableName string) string {
	return fmt.Sprintf(
		`JOIN %[1]s %[2]s ON %[2]s.%[3]s = %[4]s.%[5]s
		`,
		toBeJoinedTableName,
		nextTableName,
		query.NamingStrategy.ColumnName(nextTableName, condition.T2Field),
		previousTableName,
		query.NamingStrategy.ColumnName(previousTableName, condition.T1Field),
	)
}

// Divides a list of conditions by its type: WhereConditions and JoinConditions
func divideConditionsByType[T any](
	conditions []Condition[T],
) (thisEntityConditions []WhereCondition[T], joinConditions []Condition[T]) {
	for _, condition := range conditions {
		switch condition.Type() {
		case WhereType:
			typedCondition, ok := condition.(WhereCondition[T])
			if ok {
				// TODO ver si dejo asi
				thisEntityConditions = append(thisEntityConditions, typedCondition)
			}
		case JoinType:
			joinConditions = append(joinConditions, condition)
		}
	}

	return
}