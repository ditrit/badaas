package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ditrit/badaas/router"
)

func TestMiddlewareLogger(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()
	var actuallyRunned bool = false
	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actuallyRunned = true
	})

	router.MiddlewareLogger(nextHandler).ServeHTTP(res, req)

	if !actuallyRunned {
		t.Error("the logger middleware do not forward the request")
	}
}
