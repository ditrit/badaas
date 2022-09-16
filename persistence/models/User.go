package models

import (
	"fmt"
	"net/mail"

	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
)

// Represents a user
type User struct {
	BaseModel
	Username     string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash []byte
}

// Create a new user
func NewUser(username, email, password string) (*User, error) {
	if !validEmail(email) {
		return nil, fmt.Errorf("the email provided is not valid")
	}
	u := &User{
		Username:     username,
		Email:        email,
		PasswordHash: basicauth.SaltAndHashPassword(password),
	}
	return u, nil
}

// Validate email
func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
