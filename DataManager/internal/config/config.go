package config

import (
	"os"

	"github.com/joho/godotenv"
)

var PostgresUrl string

func Init() error {
	_ = godotenv.Load()

	// No need in container, uncomment when running locally
	//if err != nil {
	//	return errors.New("error loading .env file")
	//}

	PostgresUrl = os.Getenv("POSTGRES_URL")
	return nil
}
