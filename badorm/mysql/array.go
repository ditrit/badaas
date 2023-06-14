package mysql

import (
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
)

// Row and Array Comparisons

type ArrayExpression[T any] struct {
	Values        []T
	SQLExpression string
}

func NewArrayExpression[T any](values []T, sqlExpression string) ArrayExpression[T] {
	return ArrayExpression[T]{
		Values:        values,
		SQLExpression: sqlExpression,
	}
}

func (expr ArrayExpression[T]) ToSQL(columnName string) (string, []any) {
	placeholders := strings.Join(pie.Map(expr.Values, func(value T) string {
		return "?"
	}), ", ")

	values := pie.Map(expr.Values, func(value T) any {
		return value
	})

	return fmt.Sprintf(
		"%s %s ("+placeholders+")",
		columnName,
		expr.SQLExpression,
	), values
}

func ArrayIn[T any](values ...T) ArrayExpression[T] {
	return NewArrayExpression(values, "IN")
}

func ArrayNotIn[T any](values ...T) ArrayExpression[T] {
	return NewArrayExpression(values, "NOT IN")
}
