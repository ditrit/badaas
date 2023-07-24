package sqlservermultitype

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/multitype"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Operators
// ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16

func NotLt[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return multitype.NewValueOperator[TAttribute, TField](sql.SQLServerNotLt, field)
}

func NotGt[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return multitype.NewValueOperator[TAttribute, TField](sql.SQLServerNotGt, field)
}
