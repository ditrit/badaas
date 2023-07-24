package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/ditrit/badaas"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/testintegration"
	"github.com/ditrit/verdeter"
)

var rootCfg = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badaas",
	Short: "Backend and Distribution as a Service",
	Long:  "Badaas stands for Backend and Distribution as a Service.",
	Run:   runHTTPServer,
})

func main() {
	err := configuration.NewCommandInitializer().Init(rootCfg)
	if err != nil {
		panic(err)
	}

	rootCfg.Execute()
}

// Run the http server for badaas
func runHTTPServer(_ *cobra.Command, _ []string) {
	fx.New(
		fx.Provide(testintegration.GetModels),
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
	muxRouter *mux.Router,
	httpServerConfig configuration.HTTPServerConfiguration,
) *http.Server {
	handler := handlers.CORS(
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, "OPTIONS"}),
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
	)(muxRouter)

	srv := createServer(handler, httpServerConfig)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Sugar().Infof("Ready to serve at %s", srv.Addr)
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					log.Fatalln(err)
				}
			}()
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
