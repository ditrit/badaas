package testintegration

import (
	"github.com/stretchr/testify/suite"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/services/userservice"
)

type AuthServiceIntTestSuite struct {
	suite.Suite
	db          *badorm.DB
	userService userservice.UserService
}

func NewAuthServiceIntTestSuite(
	db *badorm.DB,
	userService userservice.UserService,
) *AuthServiceIntTestSuite {
	return &AuthServiceIntTestSuite{
		db:          db,
		userService: userService,
	}
}

func (ts *AuthServiceIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *AuthServiceIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

func (ts *AuthServiceIntTestSuite) TestGetUser() {
	email := "franco@ditrit.io"
	password := "1234"

	_, err := ts.userService.NewUser("franco", email, password)
	ts.Nil(err)

	user, err := ts.userService.GetUser(dto.UserLoginDTO{
		Email:    email,
		Password: password,
	})
	ts.Nil(err)
	ts.Equal(user.Username, "franco")
	ts.Equal(user.Email, email)
	ts.NotEqual(user.Password, password)
}
