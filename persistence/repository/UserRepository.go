package repository

import (
	"errors"

	"github.com/ditrit/badaas/persistence/models"
)

var (
	// Error "user not found"
	ErrUserNotFound = errors.New("user not found")

	// Error "user already exists"
	ErrUserAlreadyExists = errors.New("user already exists")
)

// The user registry
type UserRepository interface {
	Create(user *models.User) error
	Delete(user *models.User) error
	Save(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Search(userCriteria *models.User) ([]*models.User, error)
}
