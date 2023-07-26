package badaas

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/spf13/cobra"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/verdeter"
)

var BaDaaS = BaDaaSInitializer{}

type BaDaaSInitializer struct {
	modules []fx.Option
}

func (badaas *BaDaaSInitializer) AddModules(modules ...fx.Option) *BaDaaSInitializer {
	badaas.modules = append(badaas.modules, modules...)

	return badaas
}

func (badaas *BaDaaSInitializer) Provide(constructors ...any) *BaDaaSInitializer {
	badaas.modules = append(badaas.modules, fx.Provide(constructors...))

	return badaas
}

func (badaas BaDaaSInitializer) Init() {
	rootCfg := verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
		Use:   "badaas",
		Short: "BaDaaS",
		Run:   BaDaaS.runHTTPServer,
	})

	err := configuration.NewCommandInitializer().Init(rootCfg)
	if err != nil {
		panic(err)
	}

	rootCfg.Execute()
}

// Run the http server for badaas
func (badaas BaDaaSInitializer) runHTTPServer(cmd *cobra.Command, args []string) {
	modules := []fx.Option{
		// internal modules
		configuration.ConfigurationModule,
		router.RouterModule,
		logger.LoggerModule,
		persistence.PersistanceModule,
		services.ServicesModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		// create httpServer
		fx.Provide(newHTTPServer),
		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
	}

	fx.New(
		// add modules selected by user
		append(modules, badaas.modules...)...,
	).Run()
}
