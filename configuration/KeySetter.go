package configuration

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/models"
)

type keySetter struct{}

func newKeySetter() keySetter {
	return keySetter{}
}

// Configures the VerdeterCommand "cmd" with the information contained in "key"
func (ks keySetter) Set(cmd *verdeter.VerdeterCommand, key keyDefinition) error {
	if err := cmd.GKey(key.Name, key.ValType, "", key.Usage); err != nil {
		return err
	}

	if key.Required {
		cmd.SetRequired(key.Name)
	}

	if key.DefaultV != nil {
		cmd.SetDefault(key.Name, key.DefaultV)
	}

	if key.Validator != nil {
		cmd.AddValidator(key.Name, *key.Validator)
	}

	return nil
}

type keyDefinition struct {
	Name      string
	ValType   models.ConfigType
	Usage     string
	Required  bool
	DefaultV  any
	Validator *models.Validator
}
