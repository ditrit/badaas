package psql

import (
	"fmt"
)

// Row and Array Comparisons

type ArrayExpression[T any] struct {
	Values        []T
	SqlExpression string
}

func NewArrayExpression[T any](values []T, sqlExpression string) ArrayExpression[T] {
	return ArrayExpression[T]{
		Values:        values,
		SqlExpression: sqlExpression,
	}
}

func (expr ArrayExpression[T]) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf(
		"%s %s ?",
		columnName,
		expr.SqlExpression,
	), []any{expr.Values}
}

func ArrayIn[T any](values ...T) ArrayExpression[T] {
	return NewArrayExpression(values, "IN")
}

func ArrayNotIn[T any](values ...T) ArrayExpression[T] {
	return NewArrayExpression(values, "NOT IN")
}
