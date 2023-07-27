package configuration_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

var rootCfg = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badaas",
	Short: "Backend and Distribution as a Service",
	Run:   doNothing,
})

func doNothing(_ *cobra.Command, _ []string) {}

func TestInitCommandsInitializerConfigsAllCommandsWithoutError(t *testing.T) {
	err := configuration.NewCommandInitializer().Init(rootCfg)
	assert.Nil(t, err)
}
