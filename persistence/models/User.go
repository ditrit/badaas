package models

// Represents a user
type User struct {
	BaseModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password []byte
}
