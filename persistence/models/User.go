package models

import "github.com/ditrit/badaas/orm"

// Represents a user
type User struct {
	orm.UUIDModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`

	// password hash
	Password []byte `gorm:"not null"`
}

func UserEmailCondition(exprs ...orm.Expression[string]) orm.FieldCondition[User, string] {
	return orm.FieldCondition[User, string]{
		Expressions: exprs,
		Field:       "Email",
	}
}
