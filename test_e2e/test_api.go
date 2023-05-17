package main

import (
	"context"
	"net"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/router"
	integrationtests "github.com/ditrit/badaas/test_integration"
	"github.com/ditrit/verdeter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var rootCfg = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badaas",
	Short: "Backend and Distribution as a Service",
	Long:  "Badaas stands for Backend and Distribution as a Service.",
	Run:   runHTTPServer,
})

func main() {
	err := badaas.ConfigCommandParameters(rootCfg)
	if err != nil {
		panic(err)
	}

	rootCfg.Execute()
}

// Run the http server for badaas
func runHTTPServer(cmd *cobra.Command, args []string) {
	fx.New(
		fx.Provide(integrationtests.GetModels),
		badaas.BadaasModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		fx.Provide(NewAPIVersion),
		// add routes provided by badaas
		router.InfoRouteModule,
		router.AuthRoutesModule,
		router.EAVRoutesModule,

		fx.Invoke(badorm.AddModel[integrationtests.Company]),
		fx.Invoke(badorm.AddModel[integrationtests.Seller]),
		fx.Invoke(badorm.AddModel[integrationtests.Product]),
		router.GetCRUDRoutesModule[integrationtests.Sale](),

		// create httpServer
		fx.Provide(NewHTTPServer),

		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
	).Run()
}

func NewAPIVersion() *semver.Version {
	return semver.MustParse("0.0.0-unreleased")
}

func NewHTTPServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	router *mux.Router,
	httpServerConfig configuration.HTTPServerConfiguration,
) *http.Server {
	handler := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
			"Access-Control-Request-Headers", "Access-Control-Request-Method",
			"Connection", "Host", "Origin", "User-Agent", "Referer",
			"Cache-Control", "X-header",
		}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(0),
	)(router)

	srv := createServer(handler, httpServerConfig)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Sugar().Infof("Ready to serve at %s", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Flush the logger
			_ = logger.Sync()
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

// Create the server from the configuration holder and the http handler
func createServer(handler http.Handler, httpServerConfig configuration.HTTPServerConfiguration) *http.Server {
	timeout := httpServerConfig.GetMaxTimeout()

	return &http.Server{
		Handler: handler,
		Addr:    httpServerConfig.GetAddr(),

		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}
}
