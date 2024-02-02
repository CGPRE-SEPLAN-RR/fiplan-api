package handlers

import (
	"database/sql"

	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/database"
)

var Db *sql.DB

func StartDB() {
	Db = database.New()
}
