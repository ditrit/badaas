package commands

import (
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/ditrit/verdeter"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
	fx.New(
		// Modules
		configuration.ConfigurationModule,
		router.RouterModule,
		controllers.ControllerModule,
		logger.LoggerModule,
		persistence.PersistanceModule,

		fx.Provide(userservice.NewUserService),
		fx.Provide(sessionservice.NewSessionService),
		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		fx.Provide(NewHTTPServer),

		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
		fx.Invoke(createSuperUser),
	).Run()
	return nil
}

// The command badaas
var rootCfg = verdeter.NewVerdeterCommand(
	"badaas",
	"Backend and Distribution as a Service",
	`Badaas stands for Backend and Distribution as a Service.`,
	runHTTPServer,
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCfg.Execute()
}

func init() {
	rootCfg.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	rootCfg.SetDefault("config_path", ".")

	initServerCommands(rootCfg)
	initLoggerCommands(rootCfg)
	initDatabaseCommands(rootCfg)
	initInitialisationCommands(rootCfg)
	initSessionCommands(rootCfg)
}
