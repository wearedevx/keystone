package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	. "github.com/wearedevx/keystone/functions/ksauth"
)

func main() {
	// fmt.Println(" keystone ~ gcloud function ??,  pid=", os.Getpid(), "ppid=", os.Getppid())

	ctx := context.Background()
	// fmt.Println(" keystone ~ ksauth/main.go ~ ctx !", os.Getpid())
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", Auth); err != nil {
		log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	}
	// Use PORT environment variable, or default to 8080.
	port := "9000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	fmt.Printf("Will listen on port %s\n", port)
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
