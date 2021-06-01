package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

var (
	DB *sqlx.DB
)

func Connect() {
	log.Info("Connecting to database...")
	var err error
	DB, err = sqlx.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PWD")+"@tcp"+os.Getenv("DB_HOST")+"/"+os.Getenv("DB_NAME")+"?parseTime=true")
	if err != nil {
		log.Fatal("Connection to database failed: ", err)
		return
	}
	log.Info("Connected to database.")
}

func Close() {
	if DB == nil {
		return
	}
	log.Info("Closing DB...")
	err := DB.Close()
	if err != nil {
		log.Error("Error closing DB: ", err)
		return
	}
	log.Info("DB closed.")
}
