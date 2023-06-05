package configuration

import (
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
				Name:     DatabaseRetryKey,
				ValType:  verdeter.IsUint,
				Usage:    "The number of times badaas tries to establish a connection with the database",
				DefaultV: uint(10),
			},
			{
				Name:     DatabaseRetryDurationKey,
				ValType:  verdeter.IsUint,
				Usage:    "The duration in seconds badaas wait between connection attempts",
				DefaultV: uint(5),
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
