package mysql

import (
	"github.com/ditrit/badaas/badorm"
)

func Xor[T badorm.Model](conditions ...badorm.WhereCondition[T]) badorm.WhereCondition[T] {
	return badorm.NewConnectionCondition("XOR", conditions...)
}
