package multivalue

import (
	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/dynamic"
	"github.com/ditrit/badaas/badorm/expressions"
)

func NewMultivalueExpression[T any](sqlExpression expressions.SQLExpression, sqlConnector, sqlPrefix, sqlSuffix string, values ...T) badorm.Expression[T] {
	valuesAny := pie.Map(values, func(value T) any {
		return value
	})

	return &dynamic.MultivalueExpression[T]{
		Values:        valuesAny,
		SQLExpression: expressions.ToSQL[sqlExpression],
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
	}
}

// Equivalent to v1 < value < v2
func Between[T any](v1 T, v2 T) badorm.Expression[T] {
	return NewMultivalueExpression(expressions.Between, "AND", "", "", v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1 T, v2 T) badorm.Expression[T] {
	return NewMultivalueExpression(expressions.NotBetween, "AND", "", "", v1, v2)
}
