package sqlserverunsafe

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
	"github.com/ditrit/badaas/badorm/unsafe"
)

// Comparison Operators
// ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16

func NotLt[T any](value any) badorm.DynamicOperator[T] {
	return unsafe.NewValueOperator[T](sql.SQLServerNotLt, value)
}

func NotGt[T any](value any) badorm.DynamicOperator[T] {
	return unsafe.NewValueOperator[T](sql.SQLServerNotGt, value)
}
