package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/database"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   *sql.DB
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute,
	}

	return server
}
