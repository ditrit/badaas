package testintegration

import (
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

func TestBaDORM(t *testing.T) {
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
	tsCRUDUnsafeService *CRUDRepositoryIntTestSuite,
	db *gorm.DB,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsCRUDService)
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsCRUDUnsafeService)

	// let db cleaned
	CleanDB(db)
	shutdowner.Shutdown()
}

func NewLoggerConfiguration() configuration.LoggerConfiguration {
	viper.Set(configuration.LoggerModeKey, "dev")
	return configuration.NewLoggerConfiguration()
}

func NewGormDBConnection(logger *zap.Logger) (*gorm.DB, error) {
	dsn := "user=badaas password=badaas host=localhost port=5000 sslmode=disable dbname=badaas_db"
	// TODO codigo repetido en el ejemplo pero sin el logger
	return badorm.ConnectToDSN(logger, dsn, 10, time.Duration(5)*time.Second)
}
