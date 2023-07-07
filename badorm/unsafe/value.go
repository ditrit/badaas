package unsafe

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func NewValueOperator[T any](
	sqlOperator sql.Operator,
	value any,
) badorm.DynamicOperator[T] {
	op := badorm.NewValueOperator[T](sqlOperator, value)
	return &op
}
