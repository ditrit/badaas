package sql

type Operator uint

const (
	Eq Operator = iota
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

func (op Operator) String() string {
	return operatorToSQL[op]
}

var operatorToSQL = map[Operator]string{
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

func (op Operator) Name() string {
	return operatorToName[op]
}

var operatorToName = map[Operator]string{
	Eq:                    "Eq",
	NotEq:                 "NotEq",
	Lt:                    "Lt",
	LtOrEq:                "LtOrEq",
	Gt:                    "Gt",
	GtOrEq:                "GtOrEq",
	Between:               "Between",
	NotBetween:            "NotBetween",
	IsDistinct:            "IsDistinct",
	IsNotDistinct:         "IsNotDistinct",
	Like:                  "Like",
	Escape:                "Escape",
	MySQLIsEqual:          "mysql.IsEqual",
	MySQLRegexp:           "mysql.Regexp",
	SQLServerNotLt:        "sqlserver.NotLt",
	SQLServerNotGt:        "sqlserver.NotGt",
	PostgreSQLILike:       "psql.ILike",
	PostgreSQLSimilarTo:   "psql.SimilarTo",
	PostgreSQLPosixMatch:  "psql.PosixMatch",
	PostgreSQLPosixIMatch: "psql.PosixIMatch",
	SQLiteGlob:            "sqlite.Glob",
	ArrayIn:               "ArrayIn",
	ArrayNotIn:            "ArrayNotIn",
}
