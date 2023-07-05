package dynamic

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
func Eq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.Eq, field)
}

// NotEqualTo
func NotEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.NotEq, field)
}

// LessThan
func Lt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	// TODO invertir orden de parametros para que quede igual que en badorm/expression
	return newValueOperator(sql.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to field1 < value < field2
func Between[T any](field1 badorm.FieldIdentifier[T], field2 badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newMultivalueOperator(sql.Between, "AND", "", "", field1, field2)
}

// Equivalent to NOT (field1 < value < field2)
func NotBetween[T any](field1 badorm.FieldIdentifier[T], field2 badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newMultivalueOperator(sql.NotBetween, "AND", "", "", field1, field2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return newValueOperator(sql.IsNotDistinct, field)
}
