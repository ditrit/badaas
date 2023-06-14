package psql

import "github.com/ditrit/badaas/badorm"

// Comparison Predicates

func IsDistinct[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression(value, "IS DISTINCT FROM")
}

func IsNotDistinct[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression(value, "IS NOT DISTINCT FROM")
}
