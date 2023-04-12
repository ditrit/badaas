//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var tGlobal *testing.T

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAll(t *testing.T) {
	tGlobal = t
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

		fx.Provide(NewIntegrationTestSuite),
		fx.Provide(NewEAVServiceIntTestSuite),

		fx.Invoke(runTestSuites),
	).Run()
}

func runTestSuites(
	ts1 *EAVServiceIntTestSuite,
	shutdowner fx.Shutdowner
) {
	suite.Run(tGlobal, ts1)
	shutdowner.Shutdown()
}
