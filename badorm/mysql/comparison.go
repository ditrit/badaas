package mysql

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](value, expressions.MySQLIsEqual)
}

// Pattern Matching

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string | sql.NullString |
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](pattern string,
) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, expressions.Like)
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func RegexP[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, expressions.MySQLRegexp)
}
