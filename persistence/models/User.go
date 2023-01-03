package models

// Represents a user
type User struct {
	BaseModel
	Username       string `gorm:"unique;not null"`
	Email          string `gorm:"unique;not null"`
	OidcIdentifier string `gorm:"unique"`

	// password hash
	Password []byte `gorm:"not null"`
}

// Return the pluralized table name
//
// Satisfie the Tabler interface
func (User) TableName() string {
	return "users"
}
