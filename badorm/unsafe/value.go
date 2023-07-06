package unsafe

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func NewValueOperator[T any](
	sqlOperator sql.Operator,
	value any,
) badorm.DynamicOperator[T] {
	return &badorm.ValueOperator[T]{
		Operations: []badorm.Operation{
			{
				SQLOperator: sqlOperator,
				Value:       value,
			},
		},
		JoinNumber: badorm.UndefinedJoinNumber,
	}
}
