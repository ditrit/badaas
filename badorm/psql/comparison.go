package psql

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

// Pattern Matching

func ILike[T string | sql.NullString](pattern string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, expressions.PostgreSQLILike)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo[T string | sql.NullString](pattern string) badorm.ValueExpression[T] {
	return badorm.NewValueExpression[T](pattern, expressions.PostgreSQLSimilarTo)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, expressions.PostgreSQLPosixMatch)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, expressions.PostgreSQLPosixIMatch)
}
