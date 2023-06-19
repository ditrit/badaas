package psql

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
)

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
