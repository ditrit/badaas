package dynamic

import (
	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

func newValueExpression[T any](expression expressions.SQLExpression, field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	// TODO soportar multivalue, no todos necesariamente dinamicos
	return &badorm.ValueExpression[T]{
		SQLExpressions: []badorm.SQLExpressionAndValue{
			{
				SQL:   expressions.ToSQL[expression],
				Field: field,
			},
		},
		JoinNumber: badorm.UndefinedJoinNumber,
	}
}

func newMultivalueExpression[T any](
	sqlExpression expressions.SQLExpression,
	sqlConnector, sqlPrefix, sqlSuffix string,
	fields ...badorm.FieldIdentifier[T],
) badorm.DynamicExpression[T] {
	values := pie.Map(fields, func(field badorm.FieldIdentifier[T]) any {
		return field
	})

	return &badorm.MultivalueExpression[T]{
		Values:        values,
		SQLExpression: expressions.ToSQL[sqlExpression],
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
		JoinNumber:    badorm.UndefinedJoinNumber,
	}
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.Eq, field)
}

// NotEqualTo
func NotEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.NotEq, field)
}

// LessThan
func Lt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newValueExpression(expressions.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to field1 < value < field2
func Between[T any](field1 badorm.FieldIdentifier[T], field2 badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newMultivalueExpression(expressions.Between, "AND", "", "", field1, field2)
}

// Equivalent to NOT (field1 < value < field2)
func NotBetween[T any](field1 badorm.FieldIdentifier[T], field2 badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newMultivalueExpression(expressions.NotBetween, "AND", "", "", field1, field2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicExpression[T] {
	return newValueExpression(expressions.IsNotDistinct, field)
}
