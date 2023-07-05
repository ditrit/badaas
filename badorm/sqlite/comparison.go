package sqlite

import (
	"database/sql"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/expressions"
)

// ref: https://www.sqlite.org/lang_expr.html#like
func Glob[T string | sql.NullString](pattern string) badorm.Expression[T] {
	return badorm.NewMustBePOSIXValueExpression[T](pattern, expressions.SQLiteGlob)
}
