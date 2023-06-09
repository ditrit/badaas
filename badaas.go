package badaas

//go:generate mockery --all --keeptree

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
)

var BadaasModule = fx.Module(
	"badaas",
	configuration.ConfigurationModule,
	router.RouterModule,
	logger.LoggerModule,
	persistence.PersistanceModule,
)
