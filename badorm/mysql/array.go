package mysql

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
)

// Row and Array Comparisons

// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_in
func ArrayIn[T any](values ...T) badorm.Operator[T] {
	return badorm.NewMultivalueOperator(sql.ArrayIn, ",", "(", ")", values...)
}

func ArrayNotIn[T any](values ...T) badorm.Operator[T] {
	return badorm.NewMultivalueOperator(sql.ArrayNotIn, ",", "(", ")", values...)
}
