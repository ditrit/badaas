package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/persistence/registry"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/auth/jwtauth"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	"github.com/ditrit/badaas/services/httperrors"
)

func BasicLoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginJSONStruct dto.DTOLoginJSONStruct
	err := json.NewDecoder(r.Body).Decode(&loginJSONStruct)
	if err != nil {
		HTTPErrRequestMalformed.Write(w)
		return
	}
	user, err := registry.GetRegistry().UserRepo.GetByEmail(loginJSONStruct.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			httperrors.NewErrorNotFound(fmt.Sprintf("user %q", loginJSONStruct.Email), "the user is not in the database").Write(w)
			return
		}
		httperrors.NewInternalServerError("database error", "", err).Write(w)
		return
	}

	// Check password
	if !basicauth.CheckUserPassword(user.PasswordHash, loginJSONStruct.Password) {
		httperrors.NewUnauthorizedError("wrong password", "the provided password is incorrect").Write(w)
		return
	}

	// on valid password, generate a JWT and return it to the client
	jwtStr, err := jwtauth.CreateToken(user.Username, user.Email, user.ID)
	if err != nil {
		httperrors.NewInternalServerError("jwt error", "could not create the access token", err).Write(w)
		return
	}
	CreateAndSetCookie(w, jwtStr)
}

func CreateAndSetCookie(w http.ResponseWriter, jwtStr string) {
	accessToken := &http.Cookie{
		Name:     "access_token",
		Value:    jwtStr,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   false,
	}
	// be aware http.SetCookie may drop the cookie silently
	http.SetCookie(w, accessToken)
}

func RefreshJWT(w http.ResponseWriter, r *http.Request) {
	claims := jwtauth.JWTClaimsFromContext(r.Context())
	newAccessToken, err := jwtauth.CreateToken(claims.Username, claims.Email, uint(claims.UserID))
	if err != nil {
		// return HTTP error
		httperrors.NewInternalServerError("jwt error", "could not create the access token", err).Write(w)
		return
	}
	CreateAndSetCookie(w, newAccessToken)
}
