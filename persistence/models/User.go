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

func UserEmailCondition(email string) badorm.Condition[User] {
	return badorm.WhereCondition[User]{
		Field: "email",
		Value: email,
	}
}
