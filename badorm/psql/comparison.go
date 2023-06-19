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

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo[T string | sql.NullString](pattern string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, "SIMILAR TO")
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, "~")
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, "~*")
}
