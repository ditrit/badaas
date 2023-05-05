package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

func initDatabaseCommands(cfg *verdeter.VerdeterCommand) error {
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

	return nil
}
