package main

import (
	"github.com/gyroskan/cardinal/api"
	"github.com/gyroskan/cardinal/db"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file", err)
                log.Warn("Bonjour", err)
	}

	db.Connect()
	api.Run()
}
