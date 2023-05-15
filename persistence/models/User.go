package models

import "github.com/ditrit/badaas/badorm"

// Represents a user
type User struct {
	badorm.UUIDModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`

	// password hash
	Password []byte `gorm:"not null"`
}

// Return the pluralized table name
//
// Satisfie the Tabler interface
func (User) TableName() string {
	return "users"
}
