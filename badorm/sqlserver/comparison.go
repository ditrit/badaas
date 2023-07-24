package sqlserver

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Operators
// ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16

// EqNullable is the same as badorm.Eq but it supports value to be NULL
// ansi_nulls must be set to off to avoid the NULL = NULL: unknown problem
func EqNullable[T any](value T) badorm.Operator[T] {
	return badorm.NewValueOperator[T](sql.SQLServerEqNullable, value)
}

// NotEqNullable is the same as badorm.NotEq but it supports value to be NULL
// ansi_nulls must be set to off to avoid the NULL = NULL: unknown problem
func NotEqNullable[T any](value T) badorm.Operator[T] {
	return badorm.NewValueOperator[T](sql.SQLServerNotEqNullable, value)
}

func NotLt[T any](value T) badorm.Operator[T] {
	return badorm.NewCantBeNullValueOperator[T](sql.SQLServerNotLt, value)
}

func NotGt[T any](value T) badorm.Operator[T] {
	return badorm.NewCantBeNullValueOperator[T](sql.SQLServerNotGt, value)
}
