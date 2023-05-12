package tabler

// The interface "type" need to implement to be considered models
// TODO ver esto y los nombres de tabla que pone gorm
type Tabler interface {
	// pluralized name
	TableName() string
}
