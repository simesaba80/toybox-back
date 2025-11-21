package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DB_DSN                string
	POSTGRES_USER         string
	POSTGRES_PASSWORD     string
	POSTGRES_DB           string
	POSTGRES_HOST         string
	DISCORD_CLIENT_ID     string
	DISCORD_CLIENT_SECRET string
	TOKEN_SECRET          string
	DISCORD_GUILD_IDS     []string
	HOST_URL              string
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
	DISCORD_CLIENT_ID = os.Getenv("DISCORD_CLIENT_ID")
	DISCORD_CLIENT_SECRET = os.Getenv("DISCORD_CLIENT_SECRET")
	TOKEN_SECRET = os.Getenv("TOKEN_SECRET")
	DISCORD_GUILD_IDS = strings.Split(os.Getenv("DISCORD_GUILD_IDS"), ",")
	HOST_URL = os.Getenv("HOST_URL")
}
