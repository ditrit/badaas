package configuration

import (
	"github.com/ditrit/verdeter"
)

type CommandsInitializer interface {
	// Inits VerdeterCommand "cmd" with the all the keys that are configurable in badaas
	Init(cmd *verdeter.VerdeterCommand) error
}

type commandsInitializerImpl struct {
	KeySetter    keySetter
	Initializers []commandInitializer
}

func NewCommandInitializer() CommandsInitializer {
	return commandsInitializerImpl{
		KeySetter: newKeySetter(),
		Initializers: []commandInitializer{
			newConfigCommandInitializer(),
			newServerCommandInitializer(),
			newLoggerCommandInitializer(),
			newDatabaseCommandInitializer(),
			newInitializationCommandInitializer(),
			newSessionCommandInitializer(),
		},
	}
}

// Inits VerdeterCommand "cmd" with the all the keys that are configurable in badaas
func (ci commandsInitializerImpl) Init(cmd *verdeter.VerdeterCommand) error {
	for _, initializer := range ci.Initializers {
		if err := initializer.Init(cmd); err != nil {
			return err
		}
	}

	return nil
}

type commandInitializer struct {
	KeySetter keySetter
	Keys      []keyDefinition
}

// Inits VerdeterCommand "cmd" with the all the keys in the Keys of the initializer
func (ci commandInitializer) Init(cmd *verdeter.VerdeterCommand) error {
	for _, key := range ci.Keys {
		err := ci.KeySetter.Set(cmd, key)
		if err != nil {
			return err
		}
	}

	return nil
}
