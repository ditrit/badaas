package openid_connect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// This struct implements the functions listed in the Provider interface
type GitlabProvider struct {
	Name     string
	Config   oauth2.Config
	Verifier *oidc.IDTokenVerifier
	Provider *oidc.Provider
}

// This function creates the authentication URL to the requested provider
func (p GitlabProvider) CreateAuthURL(state string, nonce string) string {
	// access_type=offline in order to get the refresh_token
	URL := p.Config.AuthCodeURL(state, oidc.Nonce(nonce)) + "&access_type=offline"
	fmt.Println("redirectURL: " + URL + "\n")
	return URL
}

// This function exchanges the code to get the OIDC tokens from the provider
func (p GitlabProvider) GetTokens(code string) (Tokens, string, string, string) {

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
func (p GitlabProvider) RefreshTokens(refreshToken string) (Tokens, string) {

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
func (p GitlabProvider) Authenticated(rawIDToken string) AuthenticatedJson {
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
func (p GitlabProvider) RevokeToken(refreshToken string) string {

	revocation_URL, err := RevocationEndpoint(p.Provider)
	if err != nil {
		return "Failed to get the revocation_endpoint"
	}

	fmt.Println("revocation_URL: " + revocation_URL + "\n")

	client := &http.Client{}
	values := map[string]string{"token": refreshToken, "token_type_hint": "refresh_token"}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", revocation_URL, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(p.Config.ClientID, p.Config.ClientSecret)
	_, err = client.Do(req)
	if err != nil {
		return "Impossible to revoke the token"
	}

	return ""
}
