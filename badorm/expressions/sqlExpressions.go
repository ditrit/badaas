package expressions

type SQLExpression uint

const (
	Eq SQLExpression = iota
	NotEq
	Lt
	LtOrEq
	Gt
	GtOrEq
	Between
	NotBetween
	IsDistinct    // Not supported by: mysql
	IsNotDistinct // Not supported by: mysql
	Like
	Escape
	// mysql
	MySQLIsEqual
	MySQLRegexp
	// sqlserver
	SQLServerNotLt
	SQLServerNotGt
	// postgresql
	PostgreSQLILike
	PostgreSQLSimilarTo
	PostgreSQLPosixMatch
	PostgreSQLPosixIMatch
	// sqlite
	SQLiteGlob
	// shared
	ArrayIn
	ArrayNotIn
)

// alias
const (
	SQLServerEqNullable    = Eq
	SQLServerNotEqNullable = NotEq
)

var ToSQL = map[SQLExpression]string{
	Eq:                    "=",
	NotEq:                 "<>",
	Lt:                    "<",
	LtOrEq:                "<=",
	Gt:                    ">",
	GtOrEq:                ">=",
	Between:               "BETWEEN",
	NotBetween:            "NOT BETWEEN",
	IsDistinct:            "IS DISTINCT FROM",
	IsNotDistinct:         "IS NOT DISTINCT FROM",
	Like:                  "LIKE",
	Escape:                "ESCAPE",
	MySQLIsEqual:          "<=>",
	MySQLRegexp:           "REGEXP",
	SQLServerNotLt:        "!<",
	SQLServerNotGt:        "!>",
	PostgreSQLILike:       "ILIKE",
	PostgreSQLSimilarTo:   "SIMILAR TO",
	PostgreSQLPosixMatch:  "~",
	PostgreSQLPosixIMatch: "~*",
	SQLiteGlob:            "GLOB",
	ArrayIn:               "IN",
	ArrayNotIn:            "NOT IN",
}
