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
	return NewValueOperator(sql.Eq, field)
}

// NotEqualTo
func NotEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.NotEq, field)
}

// LessThan
func Lt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.Lt, field)
}

// LessThanOrEqualTo
func LtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.LtOrEq, field)
}

// GreaterThan
func Gt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.Gt, field)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.GtOrEq, field)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to field1 < value < field2
func Between[T any](field1, field2 badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewMultivalueOperator(sql.Between, sql.And, "", "", field1, field2)
}

// Equivalent to NOT (field1 < value < field2)
func NotBetween[T any](field1, field2 badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewMultivalueOperator(sql.NotBetween, sql.And, "", "", field1, field2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.IsDistinct, field)
}

// Not supported by: mysql
func IsNotDistinct[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewValueOperator(sql.IsNotDistinct, field)
}

// Row and Array Comparisons

func ArrayIn[T any](fields ...badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewMultivalueOperator(sql.ArrayIn, sql.Comma, "(", ")", fields...)
}

func ArrayNotIn[T any](fields ...badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return NewMultivalueOperator(sql.ArrayNotIn, sql.Comma, "(", ")", fields...)
}
