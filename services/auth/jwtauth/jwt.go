package jwtauth

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var secret = []byte("super secret key")

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type JWTClaims struct {
	//https://datatracker.ietf.org/doc/html/rfc7519#section-4
	jwt.RegisteredClaims

	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
}

// Create a JWT token
func CreateToken(username string, email string, id uint) (string, error) {
	// An new token instance with the Signing method
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// the expiration time. Shorter the better
	exp := time.Now().Add(time.Minute * 5)
	// now
	now := time.Now()
	token.Claims = &JWTClaims{
		jwt.RegisteredClaims{
			Issuer:    "Badaas dev instance",
			Audience:  jwt.ClaimStrings{"Badaas backend", "front ?"},
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   username,
			ID:        uuid.New().String(),
		}, username, email, int(id),
	} // Sign the token with your secret key
	str, err := token.SignedString(secret)

	if err != nil { // On error return the error
		return "", err
	} // On success, return the token string
	return str, nil
}

// Unique claim key
type claimskey int

var claimsKey claimskey

// Set KV pair in request context
func SetJWTClaimsContext(ctx context.Context, claims *JWTClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// Extract claims KV pair in request context
// Panics if the claims are not in the context
func JWTClaimsFromContext(ctx context.Context) *JWTClaims {
	claims, ok := ctx.Value(claimsKey).(*JWTClaims)
	if !ok {
		panic("could not extract claims from context")
	}
	return claims
}

// Extract claims from token string
func GetClaimsFromToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("could not get claims")
}
