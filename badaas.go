package badaas

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"go.uber.org/fx"
)

var BadaasModule = fx.Module(
	"badaas",
	configuration.ConfigurationModule,
	router.RouterModule,
	logger.LoggerModule,
	persistence.PersistanceModule,
)
