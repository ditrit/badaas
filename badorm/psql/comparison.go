package psql

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/shared"
)

func IsDistinct[T any](value T) badorm.ValueExpression[T] {
	return shared.IsDistinct(value)
}

func IsNotDistinct[T any](value T) badorm.ValueExpression[T] {
	return shared.IsNotDistinct(value)
}
