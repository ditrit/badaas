package cmd

import (
	"github.com/ditrit/verdeter"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badctl",
	Short: "the BadAas controller",
	Long:  `badctl is the command line tool that makes it possible to configure and run your BadAas applications easily.`,
})

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}
