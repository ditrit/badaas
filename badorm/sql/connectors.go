package sql

type Connector uint

const (
	And Connector = iota
	Or
	Not
	Comma
	// mysql
	MySQLXor
)

func (con Connector) String() string {
	return connectorToSQL[con]
}

var connectorToSQL = map[Connector]string{
	And:      "AND",
	Or:       "OR",
	Not:      "NOT",
	Comma:    ",",
	MySQLXor: "XOR",
}

func (con Connector) Name() string {
	return connectorToName[con]
}

var connectorToName = map[Connector]string{
	And:      "And",
	Or:       "Or",
	Not:      "Not",
	Comma:    "Comma",
	MySQLXor: "mysql.Xor",
}
