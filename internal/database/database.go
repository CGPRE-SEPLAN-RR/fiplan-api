package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/sijms/go-ora/v2"

	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal"
)

type Service interface {
	Health() map[string]string
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", username, password, host, port, database)

	db, err := sql.Open("oracle", connStr)

	if err != nil {
		log.Fatal(err)
	}

	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)

	if err != nil {
		log.Fatalf(fmt.Sprintf("O banco de dados está caído: %v", err))
	}

	return internal.BasicResponse("O banco de dados está saudável")
}

func (s *service) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args)
}

func (s *service) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args)
}
