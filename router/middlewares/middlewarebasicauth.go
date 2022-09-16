package middlewares

import (
	"errors"
	"net/http"

	"github.com/ditrit/badaas/services/auth/jwtauth"
	"github.com/ditrit/badaas/services/httperrors"
)

// The authentication middleware
func AuthenticationMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		accessTokenCookie, err := request.Cookie("access_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				httperrors.NewUnauthorizedError("Authentification Error", "No access token in cookies").Write(response)
				return
			}
			httperrors.NewInternalServerError("Authentification Error", "error while retreiving cookie", nil).Write(response)
			return
		}

		claims, err := jwtauth.GetClaimsFromToken(accessTokenCookie.Value)
		if err != nil {
			http.Error(response, err.Error(), http.StatusUnauthorized)
			return
		}
		request = request.WithContext(jwtauth.SetJWTClaimsContext(request.Context(), claims))
		next.ServeHTTP(response, request)
	})
}
