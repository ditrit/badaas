package models

var ListOfTables = []any{
	User{},
	Session{},
	EntityType{},
	Entity{},
	Value{},
	Attribute{},
	// TODO esto deberia poder ser seteado por el usuario
	Company{},
	Product{},
	Seller{},
	Sale{},
}

// The interface "type" need to implement to be considered models
type Tabler interface {
	// pluralized name
	TableName() string
}
