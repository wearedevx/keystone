package main

//go:generate go run ./generators/errors.go

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wearedevx/keystone/api/db/seed"
	"github.com/wearedevx/keystone/api/routes"
)

type baseHandler struct{}

func (h *baseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.CreateRoutes(w, r)
}

func main() {
	err := seed.SeedRoles()
	if err != nil {
		panic(err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "9001"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	server := http.Server{
		Addr:           ":" + port,
		Handler:        new(baseHandler),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatalf("Api main: %v\n", server.ListenAndServe())
}
