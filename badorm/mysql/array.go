package mysql

import "github.com/ditrit/badaas/badorm"

// Row and Array Comparisons

// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_in
func ArrayIn[T any](values ...T) badorm.MultivalueExpression[T] {
	return badorm.NewMultivalueExpression("IN", ",", "(", ")", values...)
}

func ArrayNotIn[T any](values ...T) badorm.MultivalueExpression[T] {
	return badorm.NewMultivalueExpression("NOT IN", ",", "(", ")", values...)
}
