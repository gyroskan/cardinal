package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

var (
	DB *sqlx.DB
)

func Connect() {
	connect(os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
}

func connect(db_user string , db_pwd string , db_host string, db_name string) {
	log.Info("Connecting to database...")
	var err error
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", db_user, db_pwd, db_host, db_name)
	log.Debug("Connecting with string: ", dbString)
	DB, err = sqlx.Connect("mysql", dbString)
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
