package openid_connect

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

/* This function fetches the provider well-known/openid-configuration endpoint
to get the revocation_endpoint URL */
func RevocationEndpoint(p *oidc.Provider) (string, error) {
	claims := struct {
		RevocationEndpoint string `json:"revocation_endpoint"`
	}{}
	if err := p.Claims(&claims); err != nil {
		return "", errors.Wrap(err, "Error unmarshalling provider doc into struct")
	}
	if claims.RevocationEndpoint == "" {
		return "", errors.New("Provider doesn't have a revocation_endpoint")
	}
	return claims.RevocationEndpoint, nil
}

func DoRevokeToken(ctx context.Context, revocationEndpoint string, token, tokenType string) error {
	// Verify revocation_endpoint has https url
	if !strings.HasPrefix(revocationEndpoint, "https") {
		return errors.New(fmt.Sprintf("Revocation endpoint (%v) MUST use https", revocationEndpoint))
	}
	values := url.Values{}
	values.Set("token", token)
	values.Set("token_type_hint", tokenType)
	req, err := http.NewRequest(http.MethodPost, revocationEndpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// We only support basic auth now, we may need to support other methods in the future
	// See: https://github.com/golang/oauth2/blob/bf48bf16ab8d622ce64ec6ce98d2c98f916b6303/internal/token.go#L204-L215
	// req.SetBasicAuth(clientID, clientSecret)

	resp, err := doRequest(ctx, req)
	if err != nil {
		return errors.Wrap(err, "Error contacting revocation endpoint")
	}
	if code := resp.StatusCode; code != 200 {
		// Read body to include in error for debugging purposes.
		// According to RFC6749 (https://tools.ietf.org/html/rfc6749#section-5.2)
		// the body should be in JSON, if we want to parse it in the future.
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("Revocation endpoint returned code %v, failed to read body: %v", code, err))
		}
		return errors.New(fmt.Sprintf("Revocation endpoint returned code %v, server returned: %v", code, body))
	}
	return nil
}

func doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := http.DefaultClient
	if c, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = c
	}
	// TODO: Consider retrying the request if response code is 503
	// See: https://tools.ietf.org/html/rfc7009#section-2.2.1
	return client.Do(req.WithContext(ctx))
}
