package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func initAuthenticationCommands(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.AuthTypeKey, verdeter.IsInt, "", "The type of authentication we want badaas to use.")
	cfg.SetDefault(configuration.AuthTypeKey, string(configuration.AuthTypeNone))
	cfg.AddValidator(configuration.AuthTypeKey, validators.AuthorizedValues("authorized values",
		string(configuration.AuthTypePlain),
		string(configuration.AuthTypeOIDC),
	),
	)
}
