package main

import (
	"log"

	"github.com/Daniel-Q-Reis/ia_social_media/config"
	"github.com/Daniel-Q-Reis/ia_social_media/internal/app"
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
