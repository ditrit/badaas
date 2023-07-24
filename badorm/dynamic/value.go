package dynamic

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func NewValueOperator[T any](
	sqlOperator sql.Operator,
	field badorm.FieldIdentifier[T],
) badorm.DynamicOperator[T] {
	op := badorm.NewValueOperator[T](sqlOperator, field)
	return &op
}
