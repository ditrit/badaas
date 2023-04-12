package integration_test

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	logger *zap.Logger
	db     *gorm.DB
}

func NewIntegrationTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		logger: logger,
		db:     db,
	}
}

// func (ts *IntegrationTestSuite) SetupSuite() {
// }
