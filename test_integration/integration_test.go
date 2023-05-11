package integrationtests

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/verdeter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
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
		// Modules
		badaas.BadaasModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		services.EAVServiceModule,
		fx.Provide(NewEAVServiceIntTestSuite),

		fx.Provide(persistence.AddModel[Company]),
		fx.Provide(persistence.AddModel[Seller]),
		services.GetCRUDServiceModule[Product](),
		services.GetCRUDServiceModule[Sale](),
		fx.Provide(NewCRUDServiceIntTestSuite),

		fx.Invoke(runTestSuites),
	).Run()
}

func runTestSuites(
	ts1 *EAVServiceIntTestSuite,
	ts2 *CRUDServiceIntTestSuite,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, ts1)
	suite.Run(tGlobal, ts2)
	shutdowner.Shutdown()
}
