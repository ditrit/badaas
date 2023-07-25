package testintegration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/logger"
	"github.com/ditrit/badaas/badorm/logger/gormzap"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/testintegration/models"
)

const dbTypeEnvKey = "DB"

const (
	username = "badaas"
	password = "badaas_password2023"
	host     = "localhost"
	port     = 5000
	sslMode  = "disable"
	dbName   = "badaas_db"
)

func TestBaDORM(t *testing.T) {
	tGlobal = t

	fx.New(
		fx.Provide(NewGormDBConnection),
		fx.Provide(GetModels),
		badorm.BaDORMModule,

		badorm.GetCRUDServiceModule[models.Seller](),
		badorm.GetCRUDServiceModule[models.Company](),
		badorm.GetCRUDServiceModule[models.Product](),
		badorm.GetCRUDServiceModule[models.Sale](),
		badorm.GetCRUDServiceModule[models.City](),
		badorm.GetCRUDServiceModule[models.Country](),
		badorm.GetCRUDServiceModule[models.Employee](),
		badorm.GetCRUDServiceModule[models.Bicycle](),
		badorm.GetCRUDServiceModule[models.Phone](),
		badorm.GetCRUDServiceModule[models.Brand](),
		badorm.GetCRUDServiceModule[models.Child](),

		fx.Provide(NewCRUDRepositoryIntTestSuite),
		fx.Provide(NewWhereConditionsIntTestSuite),
		fx.Provide(NewJoinConditionsIntTestSuite),
		fx.Provide(NewPreloadConditionsIntTestSuite),
		fx.Provide(NewOperatorIntTestSuite),

		fx.Invoke(runBaDORMTestSuites),
	).Run()
}

func runBaDORMTestSuites(
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	tsWhereConditions *WhereConditionsIntTestSuite,
	tsJoinConditions *JoinConditionsIntTestSuite,
	tsPreloadConditions *PreloadConditionsIntTestSuite,
	tsOperators *OperatorIntTestSuite,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsWhereConditions)
	suite.Run(tGlobal, tsJoinConditions)
	suite.Run(tGlobal, tsPreloadConditions)
	suite.Run(tGlobal, tsOperators)

	shutdowner.Shutdown()
}

func NewGormDBConnection() (*gorm.DB, error) {
	switch getDBDialector() {
	case configuration.PostgreSQL:
		return badorm.ConnectToDialector(
			logger.Default.LogMode(logger.Info),
			badorm.CreatePostgreSQLDialector(host, username, password, sslMode, dbName, port),
			10, time.Duration(5)*time.Second,
		)
	case configuration.MySQL:
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		return badorm.ConnectToDialector(
			gormzap.NewDefault(zapLogger).LogMode(logger.Info),
			badorm.CreateMySQLDialector(host, username, password, dbName, port),
			10, time.Duration(5)*time.Second,
		)
	case configuration.SQLite:
		return badorm.ConnectToDialector(
			logger.Default.LogMode(logger.Info),
			badorm.CreateSQLiteDialector(host),
			10, time.Duration(5)*time.Second,
		)
	case configuration.SQLServer:
		return badorm.ConnectToDialector(
			logger.Default.LogMode(logger.Info),
			badorm.CreateSQLServerDialector(host, username, password, dbName, port),
			10, time.Duration(5)*time.Second,
		)
	default:
		return nil, fmt.Errorf("unknown db %s", getDBDialector())
	}
}

func getDBDialector() configuration.DBDialector {
	return configuration.DBDialector(os.Getenv(dbTypeEnvKey))
}
