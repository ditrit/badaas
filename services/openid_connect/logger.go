package openid_connect

import (
	"log"
	"net/http"
)

// The goal of this middleware is only to print the method used and the API endpoint asked
func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
