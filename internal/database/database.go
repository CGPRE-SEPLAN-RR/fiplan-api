package database

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sijms/go-ora/v2"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port, _  = strconv.Atoi(os.Getenv("DB_PORT"))
	host     = os.Getenv("DB_HOST")
)

func New() *sql.DB {
	connStr := go_ora.BuildUrl(host, port, database, username, password, nil)
	db, err := sql.Open("oracle", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
