package openid_connect

import (
	"github.com/ditrit/badaas/persistence/models"
)

var providerMap = make(map[string]Provider)

// This interface is a contract for the struct GoogleProvider and GitlabProvider
type Provider interface {
	CreateAuthURL(state string, nonce string) string
	GetTokens(code string) (models.Tokens, string, string, string)
	RefreshTokens(refreshToken string) (models.Tokens, string)
	Authenticated(rawIDToken string) AuthenticatedJson
	RevokeToken(refreshToken string) string
}

// This function takes the name of the provider as a parameter and returns the corresponding provider struct
func GetProvider(name string) Provider {
	provider, ok := providerMap[name]
	if ok {
		// Provider already created
		return provider
	}
	return CreateProvider(name)

}

func CreateProvider(name string) Provider {
	switch name {
	case "google":
		providerMap["google"] = createGoogleProvider()
		return providerMap["google"]
	case "gitlab":
		providerMap["gitlab"] = createGitlabProvider()
		return providerMap["gitlab"]
	default:
		return nil
	}

}
