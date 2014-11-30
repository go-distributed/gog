package main

import (
	"github.com/go-distributed/gog/config"
	"github.com/go-distributed/gog/log"
	"github.com/go-distributed/gog/rest"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Failed to parse configuration: %v\n", err)
	}

	srv := rest.NewServer(cfg)
	log.Infof("Starting server...\n")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	return
}
