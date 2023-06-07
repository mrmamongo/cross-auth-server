package main

import (
	"github.com/mrmamongo/go-auth-tg/config"
	"github.com/mrmamongo/go-auth-tg/internal/app"
	"log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
