package sqlite

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/shared"
)

// Row and Array Comparisons

func ArrayIn[T any](values ...T) badorm.ValueExpression[T] {
	return shared.ArrayIn(values...)
}

func ArrayNotIn[T any](values ...T) badorm.ValueExpression[T] {
	return shared.ArrayNotIn(values...)
}
