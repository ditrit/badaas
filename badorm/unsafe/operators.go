package unsafe

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
func Eq[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.Eq, value)
}

// NotEqualTo
func NotEq[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.NotEq, value)
}

// LessThan
func Lt[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.GtOrEq, value)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](v1, v2 any) badorm.DynamicOperator[T] {
	return NewMultivalueOperator[T](sql.Between, sql.And, "", "", v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1, v2 any) badorm.DynamicOperator[T] {
	return NewMultivalueOperator[T](sql.NotBetween, sql.And, "", "", v1, v2)
}

// Boolean Comparison Predicates

// Not supported by: mysql
func IsDistinct[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.IsDistinct, value)
}

// Not supported by: mysql
func IsNotDistinct[T any](value any) badorm.DynamicOperator[T] {
	return NewValueOperator[T](sql.IsNotDistinct, value)
}

// Row and Array Comparisons

// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_in
func ArrayIn[T any](values ...any) badorm.DynamicOperator[T] {
	return NewMultivalueOperator[T](sql.ArrayIn, sql.Comma, "(", ")", values...)
}

func ArrayNotIn[T any](values ...any) badorm.DynamicOperator[T] {
	return NewMultivalueOperator[T](sql.ArrayNotIn, sql.Comma, "(", ")", values...)
}
