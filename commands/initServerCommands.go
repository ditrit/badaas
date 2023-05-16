package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func initServerCommands(cfg *verdeter.VerdeterCommand) error {
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
