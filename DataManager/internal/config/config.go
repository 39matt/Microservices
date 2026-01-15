package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var PostgresUrl string

func Init() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file")
	}

	PostgresUrl = os.Getenv("POSTGRES_URL")
}
