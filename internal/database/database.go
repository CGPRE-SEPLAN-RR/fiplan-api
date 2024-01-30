package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() *sql.DB {
	connStr := fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s" timezone="America/Boa_Vista"`, username, password, host, port, database)

	db, err := sql.Open("godror", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
