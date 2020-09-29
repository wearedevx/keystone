package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	. "github.com/wearedevx/keystone/functions/ksapi"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", UserService); err != nil {
		log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	}
	// Use PORT environment variable, or default to 8080.
	port := "9001"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	fmt.Printf("Will listen on port %s\n", port)
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}