package models

// Describe an object type
type EntityType struct {
	BaseModel
	Name string `gorm:"unique;not null"`

	// GORM relations
	Attributes []*Attribute
}

func (et EntityType) Equal(other EntityType) bool {
	return et.ID == other.ID &&
		et.Name == other.Name
}
