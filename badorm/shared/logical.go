package shared

import "github.com/ditrit/badaas/badorm"

func Xor[T any](conditions ...badorm.WhereCondition[T]) badorm.WhereCondition[T] {
	return badorm.NewConnectionCondition("XOR", conditions...)
}
