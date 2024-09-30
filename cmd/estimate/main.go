package main

import (
	"github.com/gomarchy/estimate/internal/api"
	"github.com/gomarchy/estimate/internal/infrastructure/database"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("unable to read .env file")
	}

	database.Connect()
	database.Migrate()

	api.Start()
}
