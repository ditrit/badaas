package multitype

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

func newMultivalueOperator[TAttribute, TField any](
	sqlOperator sql.Operator,
	sqlConnector, sqlPrefix, sqlSuffix string,
	values ...any,
) badorm.DynamicOperator[TAttribute] {
	for _, value := range values {
		_, isT1 := value.(TAttribute)
		if isT1 {
			continue
		}

		_, isField := value.(badorm.FieldIdentifier[TField])
		if isField {
			invalidOperator := verifyFieldType[TAttribute, TField]()
			if invalidOperator != nil {
				return invalidOperator
			}

			continue
		}

		return badorm.NewInvalidOperator[TAttribute](ErrParamsNotValueOrField)
	}

	return &badorm.MultivalueOperator[TAttribute]{
		Values:       values,
		SQLOperator:  sqlOperator,
		SQLConnector: sqlConnector,
		SQLPrefix:    sqlPrefix,
		SQLSuffix:    sqlSuffix,
		JoinNumber:   badorm.UndefinedJoinNumber,
	}
}
