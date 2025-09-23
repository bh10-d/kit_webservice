package main

import (
	"log"
	"os"

	"go-runner/internal/runner"
)

func main() {
	// Create new runner instance

	r, err := runner.New()
	if err != nil {
		log.Fatalf("❌ Failed to create runner: %v", err)
	}

	// Start the runner
	if err := r.Start(); err != nil {
		log.Fatalf("❌ Failed to start runner: %v", err)
		os.Exit(1)
	}
}
