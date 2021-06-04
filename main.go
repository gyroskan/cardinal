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
	}

	db.Connect()
	api.InitRouter()
	api.Run()
}
