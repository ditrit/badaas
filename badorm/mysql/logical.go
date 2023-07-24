package mysql

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func Xor[T badorm.Model](conditions ...badorm.WhereCondition[T]) badorm.WhereCondition[T] {
	return badorm.NewConnectionCondition(sql.MySQLXor, conditions...)
}
