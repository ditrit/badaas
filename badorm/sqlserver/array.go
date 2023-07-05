package sqlserver

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/shared"
)

func ArrayIn[T any](values ...T) badorm.ValueOperator[T] {
	return shared.ArrayIn(values...)
}

func ArrayNotIn[T any](values ...T) badorm.ValueOperator[T] {
	return shared.ArrayNotIn(values...)
}
