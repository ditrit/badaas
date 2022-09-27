package repository

import (
	"errors"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/httperrors"
)

var (
	// Error "user not found"
	ErrUserNotFound = errors.New("user not found")

	// Error "user already exists"
	// ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserAlreadyExists httperrors.HTTPError = httperrors.NewErrorNotFound("user", "the user is not in the database")
)

// The user registry
type UserRepository interface {
	Create(user *models.User) httperrors.HTTPError
	Delete(user *models.User) httperrors.HTTPError
	Save(user *models.User) httperrors.HTTPError
	GetByID(id uint) (*models.User, httperrors.HTTPError)
	GetByEmail(email string) (*models.User, httperrors.HTTPError)
	Search(userCriteria *models.User) ([]*models.User, httperrors.HTTPError)
}
