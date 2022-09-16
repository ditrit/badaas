package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

// The goal of this middleware is only to print the method used and the API endpoint hit
func LoggerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Sugar().Debugf("[%s]%s %s", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
