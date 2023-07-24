package configuration

import (
	"fmt"

	"github.com/ditrit/badaas/configuration/defaults"
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
				Usage:    fmt.Sprintf("Maximum timeout of the http server in second (default is %ds)", defaults.ServerTimeout),
				DefaultV: defaults.ServerTimeout,
			},
			{
				Name:     ServerHostKey,
				ValType:  verdeter.IsStr,
				Usage:    fmt.Sprintf("Address to bind (default is %s)", defaults.ServerHost),
				DefaultV: defaults.ServerHost,
			},
			{
				Name:      ServerPortKey,
				ValType:   verdeter.IsInt,
				Usage:     fmt.Sprintf("Port to bind (default is %d)", defaults.ServerPort),
				DefaultV:  defaults.ServerPort,
				Validator: &validators.CheckTCPHighPort,
			},
			{
				Name:     ServerPaginationMaxElemPerPageKey,
				ValType:  verdeter.IsUint,
				Usage:    fmt.Sprintf("The max number of records returned per page (default is %d)", defaults.ServerPaginationMaxElemPerPage),
				DefaultV: defaults.ServerPaginationMaxElemPerPage,
			},
		},
	}
}
