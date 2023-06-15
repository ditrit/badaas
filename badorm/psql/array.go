package psql

import (
	"fmt"
)

// Row and Array Comparisons

type ArrayExpression[T any] struct {
	Values        []T
	SQLExpression string
}

//nolint:unused // see inside
func (expr ArrayExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func NewArrayExpression[T any](values []T, sqlExpression string) ArrayExpression[T] {
	return ArrayExpression[T]{
		Values:        values,
		SQLExpression: sqlExpression,
	}
}

func (expr ArrayExpression[T]) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf(
		"%s %s ?",
		columnName,
		expr.SQLExpression,
	), []any{expr.Values}
}

func ArrayIn[T any](values ...T) ArrayExpression[T] {
	return NewArrayExpression(values, "IN")
}

func ArrayNotIn[T any](values ...T) ArrayExpression[T] {
	return NewArrayExpression(values, "NOT IN")
}
