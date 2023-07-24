package mysqlmultitype

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/multitype"
	"github.com/ditrit/badaas/badorm/sql"
)

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[TAttribute, TField any](field badorm.FieldIdentifier[TField]) badorm.DynamicOperator[TAttribute] {
	return multitype.NewValueOperator[TAttribute, TField](sql.MySQLIsEqual, field)
}
