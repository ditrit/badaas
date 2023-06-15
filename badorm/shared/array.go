package shared

import "github.com/ditrit/badaas/badorm"

// Row and Array Comparisons

func ArrayIn[T any](values ...T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](values, "IN")
}

func ArrayNotIn[T any](values ...T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](values, "NOT IN")
}
