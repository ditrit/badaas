package models

var ListOfTables = []any{
	Session{},
	User{},
	Value{},
	Entity{},
	Attribute{},
	EntityType{},
}

// The interface "type" need to implement to be considered models
type Tabler interface {
	// pluralized name
	TableName() string
}
