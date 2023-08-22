package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
	"github.com/russellhaering/gosaml2/types"
)

func initSAMLCommands(cfg *verdeter.VerdeterCommand) {
	
	cfg.SetConfigName("config") // name of config file (without extension)
	cfg.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name

//	viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
//	viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	cfg.AddConfigPath(".")               // optionally look for config in the working directory
	err := cfg.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	
	
	
}
