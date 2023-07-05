package shared

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

// Row and Array Comparisons

func ArrayIn[T any](values ...T) badorm.ValueOperator[T] {
	return badorm.NewValueOperator[T](values, sql.ArrayIn)
}

func ArrayNotIn[T any](values ...T) badorm.ValueOperator[T] {
	return badorm.NewValueOperator[T](values, sql.ArrayNotIn)
}
