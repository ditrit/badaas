package models

// Describe an object type
type EntityType struct {
	BaseModel
	Name string `gorm:"unique;not null"`

	// GORM relations
	Attributes []*Attribute
}
