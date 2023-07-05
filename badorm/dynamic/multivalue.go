package dynamic

import (
	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func newMultivalueOperator[T any](
	sqlOperator sql.Operator,
	sqlConnector, sqlPrefix, sqlSuffix string,
	fields ...badorm.FieldIdentifier[T],
) badorm.DynamicOperator[T] {
	values := pie.Map(fields, func(field badorm.FieldIdentifier[T]) any {
		return field
	})

	return &badorm.MultivalueOperator[T]{
		Values:       values,
		SQLOperator:  sqlOperator,
		SQLConnector: sqlConnector,
		SQLPrefix:    sqlPrefix,
		SQLSuffix:    sqlSuffix,
		JoinNumber:   badorm.UndefinedJoinNumber,
	}
}
