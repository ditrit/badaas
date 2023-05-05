package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

func initInitialisationCommands(cfg *verdeter.VerdeterCommand) error {
	err := cfg.GKey(configuration.InitializationDefaultAdminPasswordKey, verdeter.IsStr, "",
		"Set the default admin password is the admin user is not created yet.")
	if err != nil {
		return err
	}
	cfg.SetDefault(configuration.InitializationDefaultAdminPasswordKey, "admin")

	return nil
}
