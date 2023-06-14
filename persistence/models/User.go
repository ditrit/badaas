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

func UserEmailCondition(exprs ...badorm.Expression[string]) badorm.FieldCondition[User, string] {
	return badorm.FieldCondition[User, string]{
		Expressions: exprs,
		Field:       "Email",
	}
}
