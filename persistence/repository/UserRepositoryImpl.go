package repository

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/httperrors"
	"gorm.io/gorm"
)

// Return a database error
func DatabaseError(msg string, err error) httperrors.HTTPError {
	return httperrors.NewInternalServerError("database error",
		msg,
		err,
	)
}

// The implementation of the [repository.UserRepository]
type UserRepositoryImpl struct {
	db *gorm.DB
}

// The constructor of the gormUserRepository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

// Create a User
func (repo *UserRepositoryImpl) Create(user *models.User) httperrors.HTTPError {
	err := repo.db.Create(user).Error
	if err != nil {
		if gormdatabase.IsDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return DatabaseError(
			fmt.Sprintf("could not create user %q", user.Username),
			err,
		)

	}
	return nil
}

// Delete a User
func (repo *UserRepositoryImpl) Delete(user *models.User) httperrors.HTTPError {
	err := repo.db.Delete(user).Error
	if err != nil {
		return DatabaseError(
			fmt.Sprintf("could not delete user %q", user.Username),
			err,
		)
	}
	return nil
}

// Save a User
func (repo *UserRepositoryImpl) Save(user *models.User) httperrors.HTTPError {
	err := repo.db.Save(user).Error
	if err != nil {
		return DatabaseError(
			fmt.Sprintf("could not save user %q", user.Username),
			err,
		)
	}
	return nil
}

// Get a User by ID
func (repo *UserRepositoryImpl) GetByID(id uint) (*models.User, httperrors.HTTPError) {
	var u models.User
	tx := repo.db.First(&u, "id = ?", id)
	if tx.Error != nil {
		return nil, DatabaseError(
			fmt.Sprintf("could not get user by id %q", id),
			tx.Error,
		)
	}
	return &u, nil
}

// Get a user by email
func (repo *UserRepositoryImpl) GetByEmail(email string) (*models.User, httperrors.HTTPError) {
	var u models.User
	tx := repo.db.First(&u, "email = ?", email)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, DatabaseError(
				fmt.Sprintf("could not get user by email %q", email),
				tx.Error,
			)
		}
	}
	return &u, nil
}

// Get user by criteria
func (repo *UserRepositoryImpl) Search(userCriteria *models.User) ([]*models.User, httperrors.HTTPError) {
	var users []*models.User
	err := repo.db.Where(userCriteria).Find(users).Error
	if err != nil {
		return nil, DatabaseError(
			fmt.Sprintf("could not get user by user criteria (%v)", userCriteria),
			err,
		)
	}
	return users, nil
}
