package openid_connect

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// This struct implements the functions listed in the Provider interface
type GoogleProvider struct {
	Name     string
	Config   oauth2.Config
	Verifier *oidc.IDTokenVerifier
	Provider *oidc.Provider
}

// This function creates the authentication URL to the requested provider
func (p GoogleProvider) CreateAuthURL(state string, nonce string) string {
	// access_type=offline in order to get the refresh_token
	URL := p.Config.AuthCodeURL(state, oidc.Nonce(nonce)) + "&access_type=offline"
	fmt.Println("redirectURL: " + URL + "\n")
	return URL
}

// This function exchanges the code to get the OIDC tokens from the provider
func (p GoogleProvider) GetTokens(code string) (Tokens, string, string, string) {

	ctx := context.Background()
	oauth2Token, err := p.Config.Exchange(ctx, code)

	fmt.Printf("oauth2Token: %+v\n\n", oauth2Token)

	if err != nil {
		return Tokens{}, "", "", "Failed to exchange tokens"
	}

	accessToken, ok := oauth2Token.Extra("access_token").(string)
	fmt.Println("accessToken: " + accessToken + "\n")
	if !ok {
		return Tokens{}, "", "", "Failed to extract the access_token"
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	fmt.Println("rawIDToken: " + rawIDToken + "\n")
	if !ok {
		return Tokens{}, "", "", "Failed to extract the id_token"
	}
	refreshToken, ok := oauth2Token.Extra("refresh_token").(string)
	fmt.Println("refreshToken: " + refreshToken)
	fmt.Println("len(refreshToken): " + strconv.Itoa(len(refreshToken)) + "\n")

	idToken, err := p.Verifier.Verify(ctx, rawIDToken)

	if err != nil {
		return Tokens{}, "", "", "Failed to verify id_token"
	}
	var IDTokenClaims *json.RawMessage = new(json.RawMessage)
	idToken.Claims(&IDTokenClaims)

	IDTokenBody := map[string]string{}
	_ = json.Unmarshal(*IDTokenClaims, &IDTokenBody)
	email := IDTokenBody["email"]
	nonce := IDTokenBody["nonce"]

	tokens := Tokens{rawIDToken, refreshToken, accessToken}

	return tokens, email, nonce, ""

}

// This function uses the refresh_token to get new OIDC tokens
func (p GoogleProvider) RefreshTokens(refreshToken string) (Tokens, string) {

	ctx := context.Background()
	token := new(oauth2.Token)
	token.RefreshToken = refreshToken
	token.Expiry = time.Now()

	ts := p.Config.TokenSource(ctx, token)

	newToken, err := ts.Token()
	if err != nil {
		return Tokens{}, "Impossible to refresh the token"
	}

	fmt.Printf("oauth2Token: %+v\n\n", newToken)

	rawIDToken, ok := newToken.Extra("id_token").(string)
	fmt.Println("rawIDToken: " + rawIDToken + "\n")
	if !ok {
		return Tokens{}, "No id_token field in oauth2 token"
	}
	accessToken, ok := newToken.Extra("access_token").(string)
	fmt.Println("accessToken: " + accessToken + "\n")
	if !ok {
		return Tokens{}, "No access_token field in oauth2 token"
	}

	tokens := Tokens{rawIDToken, refreshToken, accessToken}

	return tokens, ""
}

// This function checks validity of the id_token
func (p GoogleProvider) Authenticated(rawIDToken string) AuthenticatedJson {
	ctx := context.Background()
	_, err := p.Verifier.Verify(ctx, rawIDToken)

	authenticated := *new(string)

	if err != nil {
		authenticated = "false"
	} else {
		authenticated = "true"
	}

	authenticated_json := AuthenticatedJson{authenticated}

	return authenticated_json
}

// This function revokes the refresh_token using the revoke endpoint of the provider
func (p GoogleProvider) RevokeToken(refreshToken string) string {
	ctx := context.Background()
	revocation_URL, err := RevocationEndpoint(p.Provider)
	if err != nil {
		return "Failed to get the revocation_endpoint"
	}

	fmt.Println("revocation_URL: " + revocation_URL + "\n")

	err = DoRevokeToken(ctx, revocation_URL, refreshToken, "refresh_token")
	if err != nil {
		return "Failed to revoke token"
	}
	return ""
}
