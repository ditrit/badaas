package shared

import "github.com/ditrit/badaas/badorm"

// Comparison Predicates

func IsDistinct[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](value, "IS DISTINCT FROM")
}

func IsNotDistinct[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](value, "IS NOT DISTINCT FROM")
}
