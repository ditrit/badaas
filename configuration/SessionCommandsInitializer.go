package configuration

import (
	"github.com/ditrit/badaas/configuration/defaults"
	"github.com/ditrit/verdeter"
)

func NewSessionCommandsInitializer() CommandsInitializer {
	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:     SessionDurationKey,
				ValType:  verdeter.IsUint,
				Usage:    "The duration of a user session in seconds",
				DefaultV: defaults.SessionDuration,
			},
			{
				Name:     SessionPullIntervalKey,
				ValType:  verdeter.IsUint,
				Usage:    "The refresh interval in seconds. Badaas refresh it's internal session cache periodically",
				DefaultV: defaults.SessionPullInterval,
			},
			{
				Name:     SessionRollIntervalKey,
				ValType:  verdeter.IsUint,
				Usage:    "The interval in which the user can renew it's session by making a request",
				DefaultV: defaults.SessionRollInterval,
			},
		},
	}
}
