package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

func initOIDCCommands(cfg *verdeter.VerdeterCommand) {
	cfg.LKey(configuration.OIDCClientIDKey, verdeter.IsStr, "", "The OIDC client ID provided by the OIDC Provider")

	cfg.LKey(configuration.OIDCClientSecretIDKey, verdeter.IsStr, "", "The OIDC client secret provided by the OIDC Provider")

	cfg.LKey(configuration.OIDCIssuerKey, verdeter.IsStr, "", "The OIDC issuer URL (example: https://accounts.google.com)")

	cfg.LKey(configuration.OIDCClaimIdentifierKey, verdeter.IsStr, "",
		"The name of the unique user identifier in the claims of the ID Token returned by the OIDC Provider.")
	cfg.SetDefault(configuration.OIDCClaimIdentifierKey, "sub")

	cfg.LKey(configuration.OIDCRedirectURLKey, verdeter.IsStr, "", "The URL of the callback on the SPA")

	cfg.LKey(configuration.OIDCScopesKey, verdeter.IsStr, "", "The scopes to request to the OIDC Provider.")
}
