package mysqldynamic

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/dynamic"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[T any](field badorm.FieldIdentifier[T]) badorm.DynamicOperator[T] {
	return dynamic.NewValueOperator[T](sql.MySQLIsEqual, field)
}
