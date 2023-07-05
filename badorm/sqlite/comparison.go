package sqlite

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

// ref: https://www.sqlite.org/lang_expr.html#like
func Glob[T string | sql.NullString](pattern string) badorm.Operator[T] {
	return badorm.NewMustBePOSIXValueOperator[T](pattern, badormSQL.SQLiteGlob)
}
