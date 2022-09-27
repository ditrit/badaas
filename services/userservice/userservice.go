package userservice

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	"github.com/ditrit/badaas/validator"
)

// Create a new user
func NewUser(username, email, password string) (*models.User, error) {
	if !validator.ValidEmail(email) {
		return nil, fmt.Errorf("the email provided is not valid")
	}
	u := &models.User{
		Username: username,
		Email:    email,
		Password: basicauth.SaltAndHashPassword(password),
	}
	return u, nil
}
