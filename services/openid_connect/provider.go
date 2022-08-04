package openid_connect

import (
	"context"
	"log"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Returns the oauth2.Config, *oidc.IDTokenVerifier and *oidc.Provider for each of the OIDC provider considered, here : Google and Gitlab. The configuration (ClientID, ClientSecret and the URL of the issuer) for each provider is located in the conf.env file
func GetProviders() (oauth2.Config, *oidc.IDTokenVerifier, *oidc.Provider, oauth2.Config, *oidc.IDTokenVerifier, *oidc.Provider) {
	ctx := context.Background()

	envErr := godotenv.Load()
	if envErr != nil {
		log.Printf("Could not load conf.env variables")
		os.Exit(1)
	}

	// Google configuration
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleIssuer := os.Getenv("GOOGLE_ISSUER")

	googleProvider, err := oidc.NewProvider(ctx, googleIssuer)
	if err != nil {
		log.Fatal(err)
	}
	googleOidcConfig := &oidc.Config{
		ClientID: googleClientID,
	}
	googleVerifier := googleProvider.Verifier(googleOidcConfig)

	googleConfig := oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     googleProvider.Endpoint(),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Gitlab configuration
	gitlabClientID := os.Getenv("GITLAB_CLIENT_ID")
	gitlabClientSecret := os.Getenv("GITLAB_CLIENT_SECRET")
	gitlabIssuer := os.Getenv("GITLAB_ISSUER")

	gitlabProvider, err := oidc.NewProvider(ctx, gitlabIssuer)
	if err != nil {
		log.Fatal(err)
	}
	gitlabOidcConfig := &oidc.Config{
		ClientID: gitlabClientID,
	}
	gitlabVerifier := gitlabProvider.Verifier(gitlabOidcConfig)

	gitlabConfig := oauth2.Config{
		ClientID:     gitlabClientID,
		ClientSecret: gitlabClientSecret,
		Endpoint:     gitlabProvider.Endpoint(),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return googleConfig, googleVerifier, googleProvider, gitlabConfig, gitlabVerifier, gitlabProvider
}

var googleConfig, googleVerifier, googleProvider, gitlabConfig, gitlabVerifier, gitlabProvider = GetProviders()

// This interface is a contract for the struct GoogleProvider and GitlabProvider
type Provider interface {
	CreateAuthURL(state string, nonce string) string
	GetTokens(code string) (models.Tokens, string, string, string)
	RefreshTokens(refreshToken string) (models.Tokens, string)
	Authenticated(rawIDToken string) AuthenticatedJson
	RevokeToken(refreshToken string) string
}

// This function takes the name of the provider as a parameter and returns the corresponding provider struct
func CreateProvider(name string) Provider {
	if name == "google" {
		return GoogleProvider{"google", googleConfig, googleVerifier, googleProvider}
	} else if name == "gitlab" {
		return GitlabProvider{"gitlab", gitlabConfig, gitlabVerifier, gitlabProvider}
	} else {
		return nil
	}
}
