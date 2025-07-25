package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Squishy ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Server Shutdown:", err)
	}

	log.Println("Squishy exiting")
}
