package cmd

import (
	"github.com/ditrit/verdeter"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badgen",
	Short: "Code generator for Badaas",
	Long:  `TODO`,
})

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}
