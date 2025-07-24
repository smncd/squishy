package main

import (
	"log"
	"net/http"

	"gitlab.com/smncd/squishy/internal/filesystem"
	"gitlab.com/smncd/squishy/internal/server"
)

func main() {
	s := &filesystem.SquishyFile{
		FilePath: "squishy.yaml",
	}

	err := s.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server := server.New(s)

	err = filesystem.LoadVersion()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Starting Squishy v%s...", filesystem.Version)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Squishy stopped.")
}
