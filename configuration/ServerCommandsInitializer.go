package configuration

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func NewServerCommandsInitializer() CommandsInitializer {
	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:     ServerTimeoutKey,
				ValType:  verdeter.IsInt,
				Usage:    "Maximum timeout of the http server in second (default is 15s)",
				DefaultV: 15,
			},
			{
				Name:     ServerHostKey,
				ValType:  verdeter.IsStr,
				Usage:    "Address to bind (default is 0.0.0.0)",
				DefaultV: "0.0.0.0",
			},
			{
				Name:      ServerPortKey,
				ValType:   verdeter.IsInt,
				Usage:     "Port to bind (default is 8000)",
				DefaultV:  8000,
				Validator: &validators.CheckTCPHighPort,
			},
			{
				Name:     ServerPaginationMaxElemPerPage,
				ValType:  verdeter.IsUint,
				Usage:    "The max number of records returned per page",
				DefaultV: uint(100),
			},
			{
				Name:     ServerExampleKey,
				ValType:  verdeter.IsStr,
				Usage:    "Example server to exec (birds | posts)",
				DefaultV: "",
			},
		},
	}
}
