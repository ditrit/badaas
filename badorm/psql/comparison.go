package psql

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

// Pattern Matching

func ILike[T string | sql.NullString](pattern string) badorm.ValueOperator[T] {
	return badorm.NewValueOperator[T](badormSQL.PostgreSQLILike, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo[T string | sql.NullString](pattern string) badorm.ValueOperator[T] {
	return badorm.NewValueOperator[T](badormSQL.PostgreSQLSimilarTo, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch[T string | sql.NullString](pattern string) badorm.Operator[T] {
	return badorm.NewMustBePOSIXValueOperator[T](badormSQL.PostgreSQLPosixMatch, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch[T string | sql.NullString](pattern string) badorm.Operator[T] {
	return badorm.NewMustBePOSIXValueOperator[T](badormSQL.PostgreSQLPosixIMatch, pattern)
}
