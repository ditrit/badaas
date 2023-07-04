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

func UserEmailCondition(expr badorm.Expression[string]) badorm.WhereCondition[User] {
	return badorm.FieldCondition[User, string]{
		Expression: expr,
		FieldIdentifier: badorm.FieldIdentifier[string]{
			Field: "Email",
		},
	}
}
