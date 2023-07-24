package unsafe

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func NewMultivalueOperator[T any](
	sqlOperator sql.Operator,
	sqlConnector sql.Connector,
	sqlPrefix, sqlSuffix string,
	values ...any,
) badorm.DynamicOperator[T] {
	return &badorm.MultivalueOperator[T]{
		Values:       values,
		SQLOperator:  sqlOperator,
		SQLConnector: sqlConnector,
		SQLPrefix:    sqlPrefix,
		SQLSuffix:    sqlSuffix,
	}
}
