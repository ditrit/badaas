package badaas

import (
	"context"
	"net"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/ditrit/badaas/configuration"
)

func newHTTPServer(
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
