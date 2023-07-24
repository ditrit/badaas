package sqlserverdynamic

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/dynamic"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Operators
// ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16

func NotLt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return dynamic.NewValueOperator[T](sql.SQLServerNotLt, field)
}

func NotGt[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return dynamic.NewValueOperator[T](sql.SQLServerNotGt, field)
}
