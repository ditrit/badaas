package registry

// Describe a type of data storage
type Datastore int

const (
	_ = iota

	// A Datastore using gorm
	//
	// please see gorm.io
	GormDatastore
)
