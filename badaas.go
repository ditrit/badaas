//Package main :
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ditrit/badaas/router"
	"github.com/gorilla/handlers"
)

// Badaas application, run a http-server on 8000.
func main() {
	router := router.SetupRouter()

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Request-Headers", "Access-Control-Request-Method", "Connection", "Host", "Origin", "User-Agent", "Referer", "Cache-Control", "X-header"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(0),
	)(router)

	srv := &http.Server{
		Handler: cors,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Listening on : " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
