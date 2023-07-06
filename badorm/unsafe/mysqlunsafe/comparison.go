package mysqlunsafe

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/sql"
	"github.com/ditrit/badaas/badorm/unsafe"
)

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[T any](value any) badorm.DynamicOperator[T] {
	return unsafe.NewValueOperator[T](sql.MySQLIsEqual, value)
}
