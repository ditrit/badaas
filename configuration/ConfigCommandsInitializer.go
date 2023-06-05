package configuration

import "github.com/ditrit/verdeter"

func NewConfigCommandsInitializer() CommandsInitializer {
	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:     "config_path",
				ValType:  verdeter.IsStr,
				Usage:    "Path to the config file/directory",
				DefaultV: ".",
			},
		},
	}
}
