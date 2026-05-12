package main

import (
	"log"

	"musicapp/backend/internal/config"
	"musicapp/backend/internal/database"
	"musicapp/backend/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Printf("connected to database at %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	r := router.New(db, cfg)

	log.Printf("server starting on port %s (environment: %s)", cfg.ServerPort, cfg.Environment)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
