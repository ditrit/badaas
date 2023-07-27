package configuration

import "github.com/ditrit/verdeter"

// Definition of config configuration keys
func newConfigCommandInitializer() commandInitializer {
	return commandInitializer{
		KeySetter: newKeySetter(),
		Keys: []keyDefinition{
			{
				Name:     "config_path",
				ValType:  verdeter.IsStr,
				Usage:    "Path to the config file/directory",
				DefaultV: ".",
			},
		},
	}
}
