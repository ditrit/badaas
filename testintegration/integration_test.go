package testintegration

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/badaas/testintegration/models"
	"github.com/ditrit/verdeter"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var tGlobal *testing.T

var testsCommand = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Run: injectDependencies,
})

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAll(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	viper.Set("config_path", path.Join(basePath, "int_test_config.yml"))
	err := badaas.ConfigCommandParameters(testsCommand)
	if err != nil {
		panic(err)
	}

	tGlobal = t

	testsCommand.Execute()
}

func injectDependencies(cmd *cobra.Command, args []string) {
	fx.New(
		fx.Provide(GetModels),
		badaas.BadaasModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		services.EAVServiceModule,
		fx.Provide(NewEAVServiceIntTestSuite),

		badorm.GetCRUDServiceModule[models.Seller, uuid.UUID](),
		badorm.GetCRUDServiceModule[models.Product, uuid.UUID](),
		badorm.GetCRUDServiceModule[models.Sale, uuid.UUID](),
		badorm.GetCRUDServiceModule[models.City, uuid.UUID](),
		badorm.GetCRUDServiceModule[models.Country, uuid.UUID](),
		badorm.GetCRUDServiceModule[models.Employee, uuid.UUID](),
		badorm.GetCRUDServiceModule[models.Bicycle, uuid.UUID](),

		badorm.GetCRUDUnsafeServiceModule[models.Company, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Seller, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Product, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Sale, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.City, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Country, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Employee, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Person, uuid.UUID](),
		badorm.GetCRUDUnsafeServiceModule[models.Bicycle, uuid.UUID](),

		fx.Provide(NewCRUDServiceIntTestSuite),
		fx.Provide(NewCRUDUnsafeServiceIntTestSuite),
		fx.Provide(NewCRUDRepositoryIntTestSuite),

		fx.Invoke(runTestSuites),
	).Run()
}

func runTestSuites(
	tsEAVService *EAVServiceIntTestSuite,
	tsCRUDService *CRUDServiceIntTestSuite,
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	tsCRUDUnsafeService *CRUDRepositoryIntTestSuite,
	db *gorm.DB,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsEAVService)
	suite.Run(tGlobal, tsCRUDService)
	suite.Run(tGlobal, tsCRUDRepository)
	suite.Run(tGlobal, tsCRUDUnsafeService)

	// let db cleaned
	CleanDB(db)
	shutdowner.Shutdown()
}
