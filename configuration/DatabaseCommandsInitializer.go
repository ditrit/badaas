package configuration

import (
	"fmt"

	"github.com/ditrit/badaas/configuration/defaults"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func NewDatabaseCommandsInitializer() CommandsInitializer {
	dialectorValidator := validators.AuthorizedValues(DBDialectors...)

	return CommandsKeyInitializer{
		KeySetter: NewKeySetter(),
		Keys: []Key{
			{
				Name:    DatabasePortKey,
				ValType: verdeter.IsInt,
				Usage:   "The port of the database server",
			},
			{
				Name:    DatabaseHostKey,
				ValType: verdeter.IsStr,
				Usage:   "The host of the database server",
			},
			{
				Name:    DatabaseNameKey,
				ValType: verdeter.IsStr,
				Usage:   "The name of the database to use",
			},
			{
				Name:    DatabaseUsernameKey,
				ValType: verdeter.IsStr,
				Usage:   "The username of the account on the database server",
			},
			{
				Name:    DatabasePasswordKey,
				ValType: verdeter.IsStr,
				Usage:   "The password of the account one the database server",
			},
			{
				Name:    DatabaseSslmodeKey,
				ValType: verdeter.IsStr,
				Usage:   "The sslmode to use when connecting to the database server",
			},
			{
				Name:    DatabaseRetryKey,
				ValType: verdeter.IsUint,
				Usage: fmt.Sprintf(
					"The number of times badaas tries to establish a connection with the database (default is %d)",
					defaults.DatabaseRetryTimes,
				),
				DefaultV: defaults.DatabaseRetryTimes,
			},
			{
				Name:    DatabaseRetryDurationKey,
				ValType: verdeter.IsUint,
				Usage: fmt.Sprintf(
					"The duration in seconds badaas wait between connection attempts (default is %ds)",
					defaults.DatabaseRetryDuration,
				),
				DefaultV: defaults.DatabaseRetryDuration,
			},
			{
				Name:      DatabaseDialectorKey,
				ValType:   verdeter.IsStr,
				Usage:     "The dialector to use to connect with the database server",
				Validator: &dialectorValidator,
			},
		},
	}
}
