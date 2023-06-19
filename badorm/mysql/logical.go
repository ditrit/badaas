package mysql

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/shared"
)

func Xor[T any](conditions ...badorm.WhereCondition[T]) badorm.WhereCondition[T] {
	return shared.Xor(conditions...)
}
