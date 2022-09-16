package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestMiddlewareLogger(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/whatever",
		},
	}
	res := httptest.NewRecorder()
	var actuallyRunned bool = false
	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actuallyRunned = true
	})

	LoggerMW(nextHandler).ServeHTTP(res, req)

	if !actuallyRunned {
		t.Error("the logger middleware did not forward the request")
	}
}
