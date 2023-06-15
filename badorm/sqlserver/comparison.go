package sqlserver

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/shared"
)

// Comparison Operators
// ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16

// EqNullable is the same as badorm.Eq but it supports value to be NULL
// ansi_nulls must be set to off to avoid the NULL = NULL: unknown problem
func EqNullable[T any](value T) badorm.Expression[T] {
	return badorm.NewValueExpression[T](value, "=")
}

// NotEqNullable is the same as badorm.NotEq but it supports value to be NULL
// ansi_nulls must be set to off to avoid the NULL = NULL: unknown problem
func NotEqNullable[T any](value T) badorm.Expression[T] {
	return badorm.NewValueExpression[T](value, "<>")
}

func NotLt[T any](value T) badorm.Expression[T] {
	return badorm.NewCantBeNullValueExpression[T](value, "!<")
}

func NotGt[T any](value T) badorm.Expression[T] {
	return badorm.NewCantBeNullValueExpression[T](value, "!>")
}

// https://learn.microsoft.com/en-us/sql/t-sql/queries/is-distinct-from-transact-sql?view=sql-server-ver16
func IsDistinct[T any](value T) badorm.ValueExpression[T] {
	return shared.IsDistinct(value)
}

func IsNotDistinct[T any](value T) badorm.ValueExpression[T] {
	return shared.IsNotDistinct(value)
}
