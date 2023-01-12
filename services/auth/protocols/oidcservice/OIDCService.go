package oidcservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/ditrit/badaas/configuration"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	// ErrFailedToExchangeAuthorizationCode is returned when the oauth2Config
	// fail to exchange the tokens against the authorization code.
	ErrFailedToExchangeAuthorizationCode = errors.New("failed to exchange authorization code")

	// ErrClaimIdentifierNotFoundInIDTokenBody is returned when the OIDC claim is not found in the ID token.
	ErrClaimIdentifierNotFoundInIDTokenBody = errors.New("claim identifier not found in ID Token body")
)

// OIDCService integrate an OIDC relying party for badaas.
type OIDCService interface {
	BuildRedirectURL(state string) string
	ExchangeAuthorizationCode(code string) (userIdentifier string, err error)
}

// Check interface compliance.
var _ OIDCService = (*oidcService)(nil)

// oidcService is a concrete implementation of OIDCService.
type oidcService struct {
	logger *zap.Logger
	// oauth2Config oauth2
	oauth2Config oauth2.Config
	// The provider OIDC
	provider *oidc.Provider

	// configuration
	oidConfiguration configuration.OIDCConfiguration
}

// BuildRedirectURL return an authentication url using the provided state
func (oidcService *oidcService) BuildRedirectURL(state string) string {
	// access_type=offline in order to get the refresh_token
	return oidcService.oauth2Config.AuthCodeURL(state)
}

// ExchangeAuthorizationCode return the identifying claim value.
func (oidcService *oidcService) ExchangeAuthorizationCode(code string) (string, error) {
	ctx := context.Background()
	oauth2Token, err := oidcService.oauth2Config.Exchange(ctx, code)

	if err != nil {
		return "", ErrFailedToExchangeAuthorizationCode
	}

	userInfo, err := oidcService.provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return "", fmt.Errorf("the UserInfo endpoint return an error: %s", err.Error())
	}

	IDTokenBody := map[string]string{}
	userInfo.Claims(&IDTokenBody)
	userIdentifier, ok := IDTokenBody[oidcService.oidConfiguration.GetClaimIdentifier()]
	if !ok {
		return "", ErrClaimIdentifierNotFoundInIDTokenBody
	}
	return userIdentifier, nil
}

// NewOIDCService is the constructor for the OIDC Service
func NewOIDCService(logger *zap.Logger, oidConfiguration configuration.OIDCConfiguration) (OIDCService, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, oidConfiguration.GetIssuer())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize oidc provider with issuer %q, error=%q",
			oidConfiguration.GetIssuer(), err.Error())
	}

	config := oauth2.Config{
		ClientID:     oidConfiguration.GetClientID(),
		ClientSecret: oidConfiguration.GetClientSecret(),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  oidConfiguration.GetRedirectURL(),
		Scopes:       append(oidConfiguration.GetScopes(), oidc.ScopeOpenID),
	}
	return &oidcService{logger, config, provider, oidConfiguration}, nil
}
