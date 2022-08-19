package openid_connect

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ditrit/badaas/persistence/models"
	"golang.org/x/oauth2"
)

type GenericOIDCProvider struct {
	Name     string
	Config   oauth2.Config
	Verifier *oidc.IDTokenVerifier
	Provider *oidc.Provider
}

func NewOIDCProvider(name, clientId, clientSecret, issuer string, scopes []string, callbackURI string) (OIDCProvider, error) {
	ctx := context.Background()

	Provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal(err)
	}
	OidcConfig := &oidc.Config{
		ClientID: clientId,
	}
	Verifier := Provider.Verifier(OidcConfig)

	Config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     Provider.Endpoint(),
		RedirectURL:  callbackURI,
		Scopes:       append(scopes, oidc.ScopeOpenID),
	}
	return GenericOIDCProvider{name, Config, Verifier, Provider}, nil
}

// This function creates the authentication URL to the requested provider
func (p GenericOIDCProvider) CreateAuthURL(state string, nonce string) string {
	// access_type=offline in order to get the refresh_token
	URL := p.Config.AuthCodeURL(state, oidc.Nonce(nonce)) + "&access_type=offline"
	log.Println("redirectURL: " + URL + "\n")
	return URL
}

// This function exchanges the code to get the OIDC tokens from the provider
func (p GenericOIDCProvider) GetTokens(code string) (models.Tokens, string, string, string) {

	ctx := context.Background()
	oauth2Token, err := p.Config.Exchange(ctx, code)

	log.Printf("oauth2Token: %+v\n\n", oauth2Token)

	if err != nil {
		return models.Tokens{}, "", "", "Failed to exchange tokens"
	}

	accessToken, ok := oauth2Token.Extra("access_token").(string)
	log.Printf("accessToken: %s\n", accessToken)
	if !ok {
		return models.Tokens{}, "", "", "Failed to extract the access_token"
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	log.Println("rawIDToken: " + rawIDToken + "\n")
	if !ok {
		return models.Tokens{}, "", "", "Failed to extract the id_token"
	}
	refreshToken, _ := oauth2Token.Extra("refresh_token").(string)
	log.Println("refreshToken: " + refreshToken)
	log.Println("len(refreshToken): " + strconv.Itoa(len(refreshToken)) + "\n")

	idToken, err := p.Verifier.Verify(ctx, rawIDToken)

	if err != nil {
		return models.Tokens{}, "", "", "Failed to verify id_token"
	}
	var IDTokenClaims *json.RawMessage = new(json.RawMessage)
	idToken.Claims(&IDTokenClaims)

	IDTokenBody := map[string]string{}
	_ = json.Unmarshal(*IDTokenClaims, &IDTokenBody)
	email := IDTokenBody["email"]
	nonce := IDTokenBody["nonce"]

	tokens := models.Tokens{rawIDToken, refreshToken, accessToken}

	return tokens, email, nonce, ""
}

// This function uses the refresh_token to get new OIDC tokens
func (p GenericOIDCProvider) RefreshTokens(refreshToken string) (models.Tokens, string) {

	ctx := context.Background()
	token := new(oauth2.Token)
	token.RefreshToken = refreshToken
	token.Expiry = time.Now()

	ts := p.Config.TokenSource(ctx, token)

	newToken, err := ts.Token()
	if err != nil {
		return models.Tokens{}, "Impossible to refresh the token"
	}

	log.Printf("oauth2Token: %+v\n\n", newToken)

	rawIDToken, ok := newToken.Extra("id_token").(string)
	log.Println("rawIDToken: " + rawIDToken + "\n")
	if !ok {
		return models.Tokens{}, "No id_token field in oauth2 token"
	}
	accessToken, ok := newToken.Extra("access_token").(string)
	log.Println("accessToken: " + accessToken + "\n")
	if !ok {
		return models.Tokens{}, "No access_token field in oauth2 token"
	}

	tokens := models.Tokens{rawIDToken, refreshToken, accessToken}

	return tokens, ""
}

// This function checks validity of the id_token
func (p GenericOIDCProvider) Authenticated(rawIDToken string) AuthenticatedJson {
	ctx := context.Background()
	_, err := p.Verifier.Verify(ctx, rawIDToken)

	if err != nil {
		return AuthenticatedJson{Value: "false"}
	}
	return AuthenticatedJson{Value: "true"}

}

// This function revokes the refresh_token using the revoke endpoint of the provider
func (p GenericOIDCProvider) RevokeToken(refreshToken string) error {
	ctx := context.Background()
	revocation_URL, err := RevocationEndpoint(p.Provider)
	if err != nil {
		return fmt.Errorf("failed to get the revocation_endpoint")
	}

	log.Println("revocation_URL: " + revocation_URL + "\n")

	err = DoRevokeToken(ctx, revocation_URL, refreshToken, "refresh_token")
	if err != nil {
		return fmt.Errorf("Failed to revoke token")
	}
	return nil
}
