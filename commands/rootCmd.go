package commands

import (
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/examples"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/resources"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/verdeter"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// Run the http server for badaas
func runHTTPServer(cmd *cobra.Command, args []string) {
	fx.New(
		// Modules
		configuration.ConfigurationModule,
		router.RouterModule,
		controllers.ControllerModule,
		logger.LoggerModule,
		persistence.PersistanceModule,
		services.ServicesModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		fx.Provide(NewHTTPServer),

		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
		fx.Invoke(createSuperUser),
		fx.Invoke(examples.StartExample),
	).Run()
}

// The command badaas
var rootCfg = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:     "badaas",
	Short:   "Backend and Distribution as a Service",
	Long:    "Badaas stands for Backend and Distribution as a Service.",
	Version: resources.Version,
	Run:     runHTTPServer,
})

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := InitCommands(rootCfg)
	if err != nil {
		panic(err)
	}

	rootCfg.Execute()
}

func InitCommands(config *verdeter.VerdeterCommand) error {
	if err := config.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory"); err != nil {
		return err
	}
	config.SetDefault("config_path", ".")

	if err := initServerCommands(config); err != nil {
		return err
	}
	if err := initLoggerCommands(config); err != nil {
		return err
	}
	if err := initDatabaseCommands(config); err != nil {
		return err
	}
	if err := initInitialisationCommands(config); err != nil {
		return err
	}
	if err := initSessionCommands(config); err != nil {
		return err
	}

	return nil
}
