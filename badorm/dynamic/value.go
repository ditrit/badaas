package dynamic

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func NewValueOperator[T any](
	sqlOperator sql.Operator,
	field badorm.FieldIdentifier[T],
) badorm.DynamicOperator[T] {
	return &badorm.ValueOperator[T]{
		Operations: []badorm.Operation{
			{
				SQLOperator: sqlOperator,
				Value:       field,
			},
		},
		JoinNumber: badorm.UndefinedJoinNumber,
	}
}
