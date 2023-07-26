package testintegration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/zap"

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
		fx.Provide(NewLogger),
		fx.Provide(NewDBConnection),
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

func NewLogger() (logger.Interface, error) {
	switch getDBDialector() {
	case configuration.PostgreSQL, configuration.SQLite, configuration.SQLServer:
		return logger.Default.ToLogMode(logger.Info), nil
	case configuration.MySQL:
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		return gormzap.NewDefault(zapLogger).ToLogMode(logger.Info), nil
	default:
		return nil, fmt.Errorf("unknown db %s", getDBDialector())
	}
}

func NewDBConnection(logger logger.Interface) (*badorm.DB, error) {
	config := badorm.Config{
		Logger:      logger,
		RetryAmount: 10,
		RetryTime:   time.Duration(5) * time.Second,
	}

	switch getDBDialector() {
	case configuration.PostgreSQL:
		return badorm.Open(
			badorm.CreatePostgreSQLDialector(host, username, password, sslMode, dbName, port),
			config,
		)
	case configuration.MySQL:
		return badorm.Open(
			badorm.CreateMySQLDialector(host, username, password, dbName, port),
			config,
		)
	case configuration.SQLite:
		return badorm.Open(
			badorm.CreateSQLiteDialector(host),
			config,
		)
	case configuration.SQLServer:
		return badorm.Open(
			badorm.CreateSQLServerDialector(host, username, password, dbName, port),
			config,
		)
	default:
		return nil, fmt.Errorf("unknown db %s", getDBDialector())
	}
}

func getDBDialector() configuration.DBDialector {
	return configuration.DBDialector(os.Getenv(dbTypeEnvKey))
}
