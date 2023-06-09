package testintegration

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/verdeter"
)

var tGlobal *testing.T

var testsCommand = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Run: injectDependencies,
})

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBaDaaS(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	viper.Set("config_path", path.Join(basePath, "int_test_config.yml"))
	viper.Set(configuration.DatabaseDialectorKey, os.Getenv(dbTypeEnvKey))

	err := configuration.NewCommandInitializer().Init(testsCommand)
	if err != nil {
		panic(err)
	}

	tGlobal = t

	testsCommand.Execute()
}

func injectDependencies(_ *cobra.Command, _ []string) {
	fx.New(
		fx.Provide(GetModels),
		badaas.BadaasModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		services.EAVServiceModule,
		fx.Provide(NewEAVServiceIntTestSuite),

		fx.Invoke(runBaDaaSTestSuites),
	).Run()
}

func runBaDaaSTestSuites(
	tsEAVService *EAVServiceIntTestSuite,
	shutdowner fx.Shutdowner,
) {
	suite.Run(tGlobal, tsEAVService)

	shutdowner.Shutdown()
}
