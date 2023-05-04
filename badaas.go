package badaas

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"go.uber.org/fx"
)

var BadaasModule = fx.Module(
	"badaas",
	configuration.ConfigurationModule,
	router.RouterModule,
	controllers.ControllerModule,
	logger.LoggerModule,
	persistence.PersistanceModule,
	services.ServicesModule,
)

func ConfigCommandParameters(command *verdeter.VerdeterCommand) {
	command.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	command.SetDefault("config_path", ".")

	configServerParameters(command)
	configLoggerParameters(command)
	configDatabaseParameters(command)
	configInitialisationParameters(command)
	configSessionParameters(command)
}

func configDatabaseParameters(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.DatabasePortKey, verdeter.IsInt, "", "The port of the database server")
	cfg.SetRequired(configuration.DatabasePortKey)

	cfg.GKey(configuration.DatabaseHostKey, verdeter.IsStr, "", "The host of the database server")
	cfg.SetRequired(configuration.DatabaseHostKey)

	cfg.GKey(configuration.DatabaseNameKey, verdeter.IsStr, "", "The name of the database to use")
	cfg.SetRequired(configuration.DatabaseNameKey)

	cfg.GKey(configuration.DatabaseUsernameKey, verdeter.IsStr, "", "The username of the account on the database server")
	cfg.SetRequired(configuration.DatabaseUsernameKey)

	cfg.GKey(configuration.DatabasePasswordKey, verdeter.IsStr, "", "The password of the account one the database server")
	cfg.SetRequired(configuration.DatabasePasswordKey)

	cfg.GKey(configuration.DatabaseSslmodeKey, verdeter.IsStr, "", "The sslmode to use when connecting to the database server")
	cfg.SetRequired(configuration.DatabaseSslmodeKey)

	cfg.GKey(configuration.DatabaseRetryKey, verdeter.IsUint, "", "The number of times badaas tries to establish a connection with the database")
	cfg.SetDefault(configuration.DatabaseRetryKey, uint(10))

	cfg.GKey(configuration.DatabaseRetryDurationKey, verdeter.IsUint, "", "The duration in seconds badaas wait between connection attempts")
	cfg.SetDefault(configuration.DatabaseRetryDurationKey, uint(5))
}

func configInitialisationParameters(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.InitializationDefaultAdminPasswordKey, verdeter.IsStr, "",
		"Set the default admin password is the admin user is not created yet.")
	cfg.SetDefault(configuration.InitializationDefaultAdminPasswordKey, "admin")
}

func configLoggerParameters(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.LoggerModeKey, verdeter.IsStr, "", "The logger mode (default to \"prod\")")
	cfg.SetDefault(configuration.LoggerModeKey, "prod")
	cfg.AddValidator(configuration.LoggerModeKey, validators.AuthorizedValues("prod", "dev"))

	cfg.GKey(configuration.LoggerRequestTemplateKey, verdeter.IsStr, "", "Template message for all request logs")
	cfg.SetDefault(configuration.LoggerRequestTemplateKey, "Receive {{method}} request on {{url}}")
}

func configServerParameters(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.ServerTimeoutKey, verdeter.IsInt, "", "Maximum timeout of the http server in second (default is 15s)")
	cfg.SetDefault(configuration.ServerTimeoutKey, 15)

	cfg.GKey(configuration.ServerHostKey, verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	cfg.SetDefault(configuration.ServerHostKey, "0.0.0.0")

	cfg.GKey(configuration.ServerPortKey, verdeter.IsInt, "p", "Port to bind (default is 8000)")
	cfg.AddValidator(configuration.ServerPortKey, validators.CheckTCPHighPort)
	cfg.SetDefault(configuration.ServerPortKey, 8000)

	cfg.GKey(configuration.ServerPaginationMaxElemPerPage, verdeter.IsUint, "", "The max number of records returned per page")
	cfg.SetDefault(configuration.ServerPaginationMaxElemPerPage, 100)

	cfg.GKey(configuration.ServerExampleKey, verdeter.IsStr, "", "Example server to exec (birds | posts)")
	cfg.SetDefault(configuration.ServerExampleKey, "")
}

func configSessionParameters(cfg *verdeter.VerdeterCommand) {
	cfg.LKey(configuration.SessionDurationKey, verdeter.IsUint, "", "The duration of a user session in seconds.")
	cfg.SetDefault(configuration.SessionDurationKey, uint(3600*4)) // 4 hours by default

	cfg.LKey(configuration.SessionPullIntervalKey,
		verdeter.IsUint, "", "The refresh interval in seconds. Badaas refresh it's internal session cache periodically.")
	cfg.SetDefault(configuration.SessionPullIntervalKey, uint(30)) // 30 seconds by default

	cfg.LKey(configuration.SessionRollIntervalKey, verdeter.IsUint, "", "The interval in which the user can renew it's session by making a request.")
	cfg.SetDefault(configuration.SessionRollIntervalKey, uint(3600)) // 1 hour by default
}
