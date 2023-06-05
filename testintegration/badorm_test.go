package testintegration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/testintegration/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const dbTypeEnvKey = "DB"

const (
	username = "badaas"
	password = "badaas"
	host     = "localhost"
	port     = 5000
	sslMode  = "disable"
	dbName   = "badaas_db"
)

func TestBaDORM(t *testing.T) {
	tGlobal = t

	fx.New(
		fx.Provide(NewLoggerConfiguration),
		logger.LoggerModule,
		fx.Provide(NewGormDBConnection),
		fx.Provide(GetModels),
		badorm.BaDORMModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		badorm.GetCRUDServiceModule[models.Seller](),
		badorm.GetCRUDServiceModule[models.Product](),
		badorm.GetCRUDServiceModule[models.Sale](),
		badorm.GetCRUDServiceModule[models.City](),
		badorm.GetCRUDServiceModule[models.Country](),
		badorm.GetCRUDServiceModule[models.Employee](),
		badorm.GetCRUDServiceModule[models.Bicycle](),

		badorm.GetCRUDUnsafeServiceModule[models.Company](),
		badorm.GetCRUDUnsafeServiceModule[models.Seller](),
		badorm.GetCRUDUnsafeServiceModule[models.Product](),
		badorm.GetCRUDUnsafeServiceModule[models.Sale](),
		badorm.GetCRUDUnsafeServiceModule[models.City](),
		badorm.GetCRUDUnsafeServiceModule[models.Country](),
		badorm.GetCRUDUnsafeServiceModule[models.Employee](),
		badorm.GetCRUDUnsafeServiceModule[models.Person](),
		badorm.GetCRUDUnsafeServiceModule[models.Bicycle](),

		fx.Provide(NewCRUDServiceIntTestSuite),
		fx.Provide(NewCRUDUnsafeServiceIntTestSuite),
		fx.Provide(NewCRUDRepositoryIntTestSuite),

		fx.Invoke(runBaDORMTestSuites),
	).Run()
}

func runBaDORMTestSuites(
	tsCRUDService *CRUDServiceIntTestSuite,
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	tsCRUDUnsafeService *CRUDUnsafeServiceIntTestSuite,
	db *gorm.DB,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsCRUDService)
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsCRUDUnsafeService)

	shutdowner.Shutdown()
}

func NewLoggerConfiguration() configuration.LoggerConfiguration {
	viper.Set(configuration.LoggerModeKey, "dev")
	return configuration.NewLoggerConfiguration()
}

func NewGormDBConnection(logger *zap.Logger) (*gorm.DB, error) {
	dbType := configuration.DBDialector(os.Getenv(dbTypeEnvKey))
	switch dbType {
	case configuration.PostgreSQL:
		return badorm.ConnectToDialector(
			logger,
			badorm.CreatePostgreSQLDialector(host, username, password, sslMode, dbName, port),
			10, time.Duration(5)*time.Second,
		)
	case configuration.MySQL:
		return badorm.ConnectToDialector(
			logger,
			badorm.CreateMySQLDialector(host, username, password, sslMode, dbName, port),
			10, time.Duration(5)*time.Second,
		)
	default:
		return nil, fmt.Errorf("unknown db %s", dbType)
	}
}
