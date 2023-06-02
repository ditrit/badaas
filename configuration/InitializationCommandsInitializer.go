package configuration

import (
	"github.com/ditrit/verdeter"
)

func NewInitializationCommandsInitializer() CommandsInitializer {
	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:     InitializationDefaultAdminPasswordKey,
				ValType:  verdeter.IsStr,
				Usage:    "Set the default admin password is the admin user is not created yet.",
				DefaultV: "admin",
			},
		},
	}
}
