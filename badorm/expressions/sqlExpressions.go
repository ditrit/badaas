package expressions

type SQLExpression uint

const (
	Eq SQLExpression = iota
	NotEq
	Lt
	LtOrEq
	Gt
	GtOrEq
	IsDistinct    // Not supported by: mysql
	IsNotDistinct // Not supported by: mysql
	// mysql
	MySQLIsEqual
	// sqlserver
	SQLServerNotLt
	SQLServerNotGt
)

// alias
const (
	SQLServerEqNullable    = Eq
	SQLServerNotEqNullable = NotEq
)

var ToSQL = map[SQLExpression]string{
	Eq:             "=",
	NotEq:          "<>",
	Lt:             "<",
	LtOrEq:         "<=",
	Gt:             ">",
	GtOrEq:         ">=",
	IsDistinct:     "IS DISTINCT FROM",
	IsNotDistinct:  "IS NOT DISTINCT FROM",
	MySQLIsEqual:   "<=>",
	SQLServerNotLt: "!<",
	SQLServerNotGt: "!>",
}
