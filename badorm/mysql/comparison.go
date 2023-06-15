package mysql

import "github.com/ditrit/badaas/badorm"

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression(value, "<=>")
}
