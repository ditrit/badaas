package jwtauth_test

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ditrit/badaas/services/auth/jwtauth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// create some fixtures
func createDummyClaims() *jwtauth.JWTClaims {
	username, email, id := "qsdfsdfsdf", "sdfsdfs@email.com", 12
	now := time.Now()
	exp := now.Add(time.Minute * 5)
	claims := &jwtauth.JWTClaims{
		jwt.RegisteredClaims{
			Issuer:    "Badaas dev instance",
			Audience:  jwt.ClaimStrings{"Badaas backend", "front ?"},
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   username,
			ID:        uuid.New().String(),
		}, username, email, int(id),
	}
	return claims
}

var regexJWT = regexp.MustCompile(`^(?:[\w-]*\.){2}[\w-]*$`)

func TestCreateToken(t *testing.T) {
	tokenStr, err := jwtauth.CreateToken("u1", "u1@email.com", 2)
	if err != nil {
		t.Errorf("should not return an error on valid args, ERROR=%s", err.Error())
	}
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		t.Errorf("the token does not respect the standard structure for a JWT. should have found 3 parts, found %d", len(parts))
	}

	if !regexJWT.Match([]byte(tokenStr)) {
		t.Errorf("the token did not match the regex")
	}
}

func TestGetClaimsFromToken(t *testing.T) {
	// old token
	tok := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJCYWRhYXMgZGV2IGluc3RhbmNlIiwic3ViIjoidTEiLCJhdWQiOlsiQmFkYWFzIGJhY2tlbmQiLCJmcm9udCA_Il0sImV4cCI6MTY2MTI1NjIxNSwibmJmIjoxNjYxMjU1OTE1LCJpYXQiOjE2NjEyNTU5MTUsImp0aSI6ImMyYzFhYWNkLTZlYjktNGM4NC1iMTQ3LTFlMmU0ZDhlM2JkZSIsInVzZXJuYW1lIjoidTEiLCJlbWFpbCI6InUxQGVtYWlsLmNvbSIsInVzZXJfaWQiOjc5MDMxNjY4MjExMzA1Njc2OX0.OkuXqFo1XjxTNHdmKPe3MjxeBDYCW1dLdwic8cG5-VM"

	_, err := jwtauth.GetClaimsFromToken(tok)
	if err == nil {
		t.Errorf("should return an error on expired token")
	}

	// New token
	tok, _ = jwtauth.CreateToken("whatever", "whatever@email.com", 2)
	claims, err := jwtauth.GetClaimsFromToken(tok)
	if err != nil {
		t.Errorf("should not return an error on valid token string")
	}
	if claims == nil {
		t.Errorf("should not return a nil pointer to JWTClaim on valid token string")
	}
}

func TestStoreClaimsInContext(t *testing.T) {
	claims := createDummyClaims()
	ctx := context.Background()

	ctx = jwtauth.SetJWTClaimsContext(ctx, claims)

	claimsExtracted := jwtauth.JWTClaimsFromContext(ctx)
	if claimsExtracted == nil {
		t.Errorf("should not return a nil pointer to JWTClaims on valid context")
	}
	if !reflect.DeepEqual(claims, claimsExtracted) {
		t.Errorf("look's like the claims extracted from the context don't exactly match the original claims")
	}
}

func TestStoreClaimsInContextPanics(t *testing.T) {
	ctx := context.Background()
	assert.Panics(t, func() { jwtauth.JWTClaimsFromContext(ctx) }, "should panic")
}
