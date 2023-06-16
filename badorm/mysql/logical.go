package mysql

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/shared"
)

func Xor[T any](exprs ...badorm.Expression[T]) badorm.ConnectionExpression[T] {
	return shared.Xor(exprs...)
}
