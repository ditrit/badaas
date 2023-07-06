package unsafe

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func NewValueOperator[T any](
	sqlOperator sql.Operator,
	value any,
) badorm.DynamicOperator[T] {
	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	// TODO soportar multivalue, no todos necesariamente dinamicos
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
