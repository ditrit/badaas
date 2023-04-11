package integration_test

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	gormDB *gorm.DB
}

func (ts *IntegrationTestSuite) SetupSuite() {
	// TODO use dependency injection
	logger, _ := zap.NewProduction()
	gormDB, err := gormdatabase.CreateDatabaseConnectionFromConfiguration(logger, configuration.NewDatabaseConfiguration())
	if err != nil {
		panic(err)
	}

	ts.gormDB = gormDB
}
