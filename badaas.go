package badaas

//go:generate mockery --all --keeptree

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"go.uber.org/fx"
)

var BadaasModule = fx.Module(
	"badaas",
	configuration.ConfigurationModule,
	router.RouterModule,
	logger.LoggerModule,
	persistence.PersistanceModule,
)

func ConfigCommandParameters(command *verdeter.VerdeterCommand) error {
	err := command.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	if err != nil {
		return err
	}
	command.SetDefault("config_path", ".")

	err = configServerParameters(command)
	if err != nil {
		return err
	}

	err = configLoggerParameters(command)
	if err != nil {
		return err
	}

	err = configDatabaseParameters(command)
	if err != nil {
		return err
	}

	err = configInitialisationParameters(command)
	if err != nil {
		return err
	}

	err = configSessionParameters(command)
	if err != nil {
		return err
	}

	return nil
}

func configDatabaseParameters(cfg *verdeter.VerdeterCommand) error {
	err := cfg.GKey(configuration.DatabasePortKey, verdeter.IsInt, "", "The port of the database server")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabasePortKey)

	err = cfg.GKey(configuration.DatabaseHostKey, verdeter.IsStr, "", "The host of the database server")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabaseHostKey)

	err = cfg.GKey(configuration.DatabaseNameKey, verdeter.IsStr, "", "The name of the database to use")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabaseNameKey)

	err = cfg.GKey(configuration.DatabaseUsernameKey, verdeter.IsStr, "", "The username of the account on the database server")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabaseUsernameKey)

	err = cfg.GKey(configuration.DatabasePasswordKey, verdeter.IsStr, "", "The password of the account one the database server")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabasePasswordKey)

	err = cfg.GKey(configuration.DatabaseSslmodeKey, verdeter.IsStr, "", "The sslmode to use when connecting to the database server")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabaseSslmodeKey)

	err = cfg.GKey(configuration.DatabaseRetryKey, verdeter.IsUint, "", "The number of times badaas tries to establish a connection with the database")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.DatabaseRetryKey, uint(10))

	err = cfg.GKey(configuration.DatabaseRetryDurationKey, verdeter.IsUint, "", "The duration in seconds badaas wait between connection attempts")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.DatabaseRetryDurationKey, uint(5))

	err = cfg.GKey(configuration.DatabaseDialectorKey, verdeter.IsStr, "", "The dialector to use to connect with the database server")
	if err != nil {
		return err
	}
	cfg.SetRequired(configuration.DatabaseDialectorKey)
	// TODO
	// cfg.AddValidator(
	// configuration.DatabaseDialectorKey,
	// validators.AuthorizedValues(configuration.DBDialectors...),
	// )

	return nil
}

func configInitialisationParameters(cfg *verdeter.VerdeterCommand) error {
	err := cfg.GKey(configuration.InitializationDefaultAdminPasswordKey, verdeter.IsStr, "",
		"Set the default admin password is the admin user is not created yet.")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.InitializationDefaultAdminPasswordKey, "admin")

	return nil
}

func configLoggerParameters(cfg *verdeter.VerdeterCommand) error {
	err := cfg.GKey(configuration.LoggerModeKey, verdeter.IsStr, "", "The logger mode (default to \"prod\")")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.LoggerModeKey, "prod")
	cfg.AddValidator(configuration.LoggerModeKey, validators.AuthorizedValues("prod", "dev"))

	err = cfg.GKey(configuration.LoggerRequestTemplateKey, verdeter.IsStr, "", "Template message for all request logs")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.LoggerRequestTemplateKey, "Receive {{method}} request on {{url}}")

	return nil
}

func configServerParameters(cfg *verdeter.VerdeterCommand) error {
	err := cfg.GKey(configuration.ServerTimeoutKey, verdeter.IsInt, "", "Maximum timeout of the http server in second (default is 15s)")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.ServerTimeoutKey, 15)

	err = cfg.GKey(configuration.ServerHostKey, verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.ServerHostKey, "0.0.0.0")

	err = cfg.GKey(configuration.ServerPortKey, verdeter.IsInt, "p", "Port to bind (default is 8000)")
	if err != nil {
		return err
	}
	cfg.AddValidator(configuration.ServerPortKey, validators.CheckTCPHighPort)
	cfg.SetDefault(configuration.ServerPortKey, 8000)

	err = cfg.GKey(configuration.ServerPaginationMaxElemPerPage, verdeter.IsUint, "", "The max number of records returned per page")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.ServerPaginationMaxElemPerPage, 100)

	err = cfg.GKey(configuration.ServerExampleKey, verdeter.IsStr, "", "Example server to exec (birds | posts)")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.ServerExampleKey, "")

	return nil
}

func configSessionParameters(cfg *verdeter.VerdeterCommand) error {
	err := cfg.LKey(configuration.SessionDurationKey, verdeter.IsUint, "", "The duration of a user session in seconds.")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.SessionDurationKey, uint(3600*4)) // 4 hours by default

	err = cfg.LKey(configuration.SessionPullIntervalKey,
		verdeter.IsUint, "", "The refresh interval in seconds. Badaas refresh it's internal session cache periodically.")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.SessionPullIntervalKey, uint(30)) // 30 seconds by default

	err = cfg.LKey(configuration.SessionRollIntervalKey, verdeter.IsUint, "", "The interval in which the user can renew it's session by making a request.")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.SessionRollIntervalKey, uint(3600)) // 1 hour by default

	return nil
}
