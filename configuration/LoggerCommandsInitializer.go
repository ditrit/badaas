package configuration

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func NewLoggerCommandsInitializer() CommandsInitializer {
	modeValidator := validators.AuthorizedValues("prod", "dev")

	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:     LoggerRequestTemplateKey,
				ValType:  verdeter.IsStr,
				Usage:    "Template message for all request logs",
				DefaultV: "Receive {{method}} request on {{url}}",
			},
			{
				Name:      LoggerModeKey,
				ValType:   verdeter.IsStr,
				Usage:     "The logger mode (default to \"prod\")",
				DefaultV:  "prod",
				Validator: &modeValidator,
			},
		},
	}
}
