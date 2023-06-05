package configuration

import (
	"github.com/ditrit/verdeter"
)

type CommandsInitializer interface {
	Init(config *verdeter.VerdeterCommand) error
}

type CommandsInitializerImpl struct {
	KeySetter    KeySetter
	Initializers []CommandsInitializer
}

func NewCommandInitializer() CommandsInitializer {
	return CommandsInitializerImpl{
		KeySetter: NewKeySetter(),
		Initializers: []CommandsInitializer{
			NewConfigCommandsInitializer(),
			NewServerCommandsInitializer(),
			NewLoggerCommandsInitializer(),
			NewDatabaseCommandsInitializer(),
			NewInitializationCommandsInitializer(),
			NewSessionCommandsInitializer(),
		},
	}
}

func (ci CommandsInitializerImpl) Init(config *verdeter.VerdeterCommand) error {
	for _, initializer := range ci.Initializers {
		if err := initializer.Init(config); err != nil {
			return err
		}
	}

	return nil
}

type CommandsKeyInitializer struct {
	KeySetter KeySetter
	Keys      []Key
}

func (ci CommandsKeyInitializer) Init(config *verdeter.VerdeterCommand) error {
	for _, key := range ci.Keys {
		err := ci.KeySetter.Set(config, key)
		if err != nil {
			return err
		}
	}

	return nil
}
