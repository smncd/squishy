package main

import (
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/server"
)

var Version string

func main() {
	s := &filesystem.SquishyFile{}

	s.SetFilePath("squishy.yaml")

	err := s.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server := server.New(s)

	log.Printf("Starting Squishy v%s...", Version)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Squishy stopped.")
}
