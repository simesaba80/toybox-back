package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_DSN            string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_HOST     string
)

// .envを呼び出します。
func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("読み込み出来ませんでした: %v", err)
	}

	DB_DSN = os.Getenv("DB_DSN")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB = os.Getenv("POSTGRES_DB")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
}
