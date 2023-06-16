package mysql

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
)

// Comparison Predicates

// preferred over eq
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_equal-to
func IsEqual[T any](value T) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](value, "<=>")
}

// Pattern Matching

// TODO codigo repetido
// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string | sql.NullString | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](pattern string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, "LIKE")
}

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func LikeEscape[T string | sql.NullString | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](pattern string, escape rune) badorm.MultiExpressionExpression[T] {
	return badorm.NewMultiExpressionExpression[T](
		badorm.SQLExpressionAndValue{
			SQLExpression: "LIKE",
			Value:         pattern,
		},
		badorm.SQLExpressionAndValue{
			SQLExpression: "ESCAPE",
			Value:         string(escape),
		},
	)
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func RegexP[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, "REGEXP")
}
