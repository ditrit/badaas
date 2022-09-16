package repository

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"gorm.io/gorm"
)

// The implementation of the [repository.UserRepository]
type gormUserRepository struct {
	db *gorm.DB
}

// The constructor of the gormUserRepository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db}
}

// Create a User
func (repo *gormUserRepository) Create(user *models.User) error {
	err := repo.db.Create(user).Error
	if err != nil {
		if gormdatabase.IsDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("could not create user")
	}
	return nil
}

// Delete a User
func (repo *gormUserRepository) Delete(user *models.User) error {
	return repo.db.Delete(user).Error
}

// Save a User
func (repo *gormUserRepository) Save(user *models.User) error {
	err := repo.db.Save(user).Error
	if err != nil {
		return fmt.Errorf("could not save user (ERROR=%q)", err.Error())
	}
	return nil
}

// Get a User by ID
func (repo *gormUserRepository) GetByID(id uint) (*models.User, error) {
	var u models.User
	tx := repo.db.First(&u, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
	}
	return &u, nil
}

// Get a user by email
func (repo *gormUserRepository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	tx := repo.db.First(&u, "email = ?", email)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
	}
	return &u, nil
}

// Get user by criteria
func (repo *gormUserRepository) Search(userCriteria *models.User) ([]*models.User, error) {
	var users []*models.User
	err := repo.db.Where(userCriteria).Find(users).Error
	if err != nil {
		return users, err
	}
	return users, nil
}
