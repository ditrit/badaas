package configuration

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/models"
)

type KeySetter interface {
	Set(config *verdeter.VerdeterCommand, key Key) error
}

type KeySetterImpl struct{}

func NewKeySetter() KeySetter {
	return KeySetterImpl{}
}

func (ks KeySetterImpl) Set(config *verdeter.VerdeterCommand, key Key) error {
	if err := config.GKey(key.Name, key.ValType, "", key.Usage); err != nil {
		return err
	}

	if key.Required {
		config.SetRequired(key.Name)
	}

	if key.DefaultV != nil {
		config.SetDefault(key.Name, key.DefaultV)
	}

	if key.Validator != nil {
		config.AddValidator(key.Name, *key.Validator)
	}

	return nil
}

type Key struct {
	Name      string
	ValType   models.ConfigType
	Usage     string
	Required  bool
	DefaultV  any
	Validator *models.Validator
}
