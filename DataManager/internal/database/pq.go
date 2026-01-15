package database

import (
	"DataManager/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitPostgres() error {
	connStr := config.PostgresUrl
	if connStr == "" {
		return fmt.Errorf("postgres connection string is empty")

	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not connect to postgres: %v", err)

	}
	defer func() {
		if err = db.Close(); err != nil {
		}
	}()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("could not ping postgres: %v", err)
	}

	log.Println("Postgres connected successfully")
	return nil
}
