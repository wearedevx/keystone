package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wearedevx/keystone/api/routes"
)

type baseHandler struct{}

func (h *baseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.CreateRoutes(w, r)
}

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "9001"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	server := http.Server{
		Addr:           ":" + port,
		Handler:        new(baseHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Will listen on port %s\n", port)
	log.Fatalf("Api main: %v\n", server.ListenAndServe())
}
