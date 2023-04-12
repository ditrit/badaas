package integration_test

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	gormDB *gorm.DB
}

func NewIntegrationTestSuite(
	db *gorm.DB,
) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		gormDB: db,
	}
}

// func (ts *IntegrationTestSuite) SetupSuite() {
// }
