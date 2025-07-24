package main

import (
	"fmt"
	"log"

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

	router := server.New(s)

	router.Run(fmt.Sprintf("%v:%v", s.Config.Host, s.Config.Port))
}
