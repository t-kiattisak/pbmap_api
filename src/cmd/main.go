package main

import (
	"pbmap_api/src/config"
	"pbmap_api/src/internal/database"
	httpHandler "pbmap_api/src/internal/handler/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	db := config.NewDatabase(cfg)

	if err := database.Migrate(db); err != nil {
		panic(err)
	}

	httpHandler.Run(cfg, db)
}
