package shared

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

// Row and Array Comparisons

func ArrayIn[T any](values ...T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](values, expressions.ArrayIn)
}

func ArrayNotIn[T any](values ...T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](values, expressions.ArrayNotIn)
}
