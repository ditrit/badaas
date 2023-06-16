package psql

import (
	"database/sql"

	"golang.org/x/text/unicode/norm"

	"github.com/ditrit/badaas/badorm"
)

// String Functions and Operators

var normalForms = map[norm.Form]string{
	norm.NFC:  "NFC",
	norm.NFD:  "NFD",
	norm.NFKC: "NFKC",
	norm.NFKD: "NFKD",
}

// TODO que pasa con otros que mapean a string por valuer?
func IsNormalized[T string | sql.NullString](expectedNorm norm.Form) badorm.PredicateExpression[T] {
	return badorm.NewPredicateExpression[T](
		"IS " + normalForms[expectedNorm] + " NORMALIZED",
	)
}

func StartsWith[T string | sql.NullString](expectedStart string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](expectedStart, "^@")
}

// Pattern Matching

func ILike[T string | sql.NullString](pattern string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, "ILIKE")
}

func ILikeEscape[T string | sql.NullString](pattern string, escape rune) badorm.MultiExpressionExpression[T] {
	// TODO aca podria hacer un .Add o algo asi para no repetir lo de arriba
	return badorm.NewMultiExpressionExpression[T](
		badorm.SQLExpressionAndValue{
			SQLExpression: "ILIKE",
			Value:         pattern,
		},
		badorm.SQLExpressionAndValue{
			SQLExpression: "ESCAPE",
			Value:         string(escape),
		},
	)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo[T string | sql.NullString](pattern string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, "SIMILAR TO")
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarToEscape[T string | sql.NullString](pattern string, escape rune) badorm.MultiExpressionExpression[T] {
	return badorm.NewMultiExpressionExpression[T](
		badorm.SQLExpressionAndValue{
			SQLExpression: "SIMILAR TO",
			Value:         pattern,
		},
		badorm.SQLExpressionAndValue{
			SQLExpression: "ESCAPE",
			Value:         string(escape),
		},
	)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, "~")
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, "~*")
}
