package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_DSN string
)

// .envを呼び出します。
func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("読み込み出来ませんでした: %v", err)
	}

	DB_DSN = os.Getenv("DB_DSN")
}
