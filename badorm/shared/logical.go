package shared

import "github.com/ditrit/badaas/badorm"

func Xor[T any](exprs ...badorm.Expression[T]) badorm.ConnectionExpression[T] {
	return badorm.NewConnectionExpression("XOR", exprs...)
}