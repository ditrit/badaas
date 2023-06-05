package integrationtests

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/services"
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
	err := configuration.NewCommandInitializer().Init(testsCommand)
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

		fx.Invoke(badorm.AddModel[Company]),
		badorm.GetCRUDServiceModule[Seller, uuid.UUID](),
		badorm.GetCRUDServiceModule[Product, uuid.UUID](),
		badorm.GetCRUDServiceModule[Sale, uuid.UUID](),
		badorm.GetCRUDServiceModule[City, uuid.UUID](),
		badorm.GetCRUDServiceModule[Country, uuid.UUID](),
		badorm.GetCRUDServiceModule[Employee, uuid.UUID](),
		fx.Invoke(badorm.AddModel[Person]),
		badorm.GetCRUDServiceModule[Bicycle, uuid.UUID](),
		fx.Provide(NewCRUDServiceIntTestSuite),
		fx.Provide(NewCRUDRepositoryIntTestSuite),

		fx.Invoke(runTestSuites),
	).Run()
}

func runTestSuites(
	tsEAV *EAVServiceIntTestSuite,
	tsCRUD *CRUDServiceIntTestSuite,
	tsCRUDRepository *CRUDRepositoryIntTestSuite,
	db *gorm.DB,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsEAV)
	suite.Run(tGlobal, tsCRUD)
	suite.Run(tGlobal, tsCRUDRepository)

	// let db cleaned
	CleanDB(db)
	shutdowner.Shutdown()
}
