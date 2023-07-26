package userservice

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	"github.com/ditrit/badaas/utils/validators"
)

// UserService provide functions related to Users
type UserService interface {
	NewUser(username, email, password string) (*models.User, error)
	GetUser(dto.UserLoginDTO) (*models.User, error)
}

var ErrWrongPassword = errors.New("password is incorrect")

// Check interface compliance
var _ UserService = (*userServiceImpl)(nil)

// The UserService concrete implementation
type userServiceImpl struct {
	userRepository badorm.CRUDRepository[models.User, badorm.UUID]
	logger         *zap.Logger
	db             *badorm.DB
}

// UserService constructor
func NewUserService(
	logger *zap.Logger,
	userRepository badorm.CRUDRepository[models.User, badorm.UUID],
	db *badorm.DB,
) UserService {
	return &userServiceImpl{
		logger:         logger,
		userRepository: userRepository,
		db:             db,
	}
}

// Create a new user
func (service *userServiceImpl) NewUser(username, email, password string) (*models.User, error) {
	sanitizedEmail, err := validators.ValidEmail(email)
	if err != nil {
		return nil, fmt.Errorf("the provided email is not valid")
	}

	u := &models.User{
		Username: username,
		Email:    sanitizedEmail,
		Password: basicauth.SaltAndHashPassword(password),
	}

	err = service.userRepository.Create(service.db.GormDB, u)
	if err != nil {
		return nil, err
	}

	service.logger.Info(
		"Successfully created a new user",
		zap.String("email", sanitizedEmail),
		zap.String("username", username),
	)

	return u, nil
}

// Get user if the email and password provided are correct, return an error if not.
func (service *userServiceImpl) GetUser(userLoginDTO dto.UserLoginDTO) (*models.User, error) {
	user, err := service.userRepository.QueryOne(
		service.db.GormDB,
		models.UserEmailCondition(badorm.Eq(userLoginDTO.Email)),
	)
	if err != nil {
		return nil, err
	}

	// Check password
	if !basicauth.CheckUserPassword(user.Password, userLoginDTO.Password) {
		return nil, ErrWrongPassword
	}

	return user, nil
}
