package mysql

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[T any](value T) badorm.ValueOperator[T] {
	return badorm.NewValueOperator[T](badormSQL.MySQLIsEqual, value)
}

// Pattern Matching

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string | sql.NullString |
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](pattern string,
) badorm.ValueOperator[T] {
	return badorm.NewValueOperator[T](badormSQL.Like, pattern)
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func RegexP[T string | sql.NullString](pattern string) badorm.Operator[T] {
	return badorm.NewMustBePOSIXValueOperator[T](badormSQL.MySQLRegexp, pattern)
}
