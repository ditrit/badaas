package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

// initialize session related config keys
func initSessionCommands(cfg *verdeter.VerdeterCommand) error {
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
