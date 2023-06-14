package mysql

import "github.com/ditrit/badaas/badorm"

// Comparison Predicates

func IsEqual[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression(value, "<=>")
}
