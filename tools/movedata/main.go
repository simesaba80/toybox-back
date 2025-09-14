package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"github.com/simesaba80/toybox-back/pkg/db"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func main() {
	// This is a placeholder main function.
	// The actual migration logic would be implemented here.
	config.LoadEnv()
	fmt.Println("Environment variables loaded.")
	//移行後のDB接続
	db.Init()
	fmt.Println("Connected to the database1.")
	//移行前のDB接続
	err := godotenv.Load()
	if err != nil {
		log.Printf("読み込み出来ませんでした: %v", err)
	}

	DB_DSN_BACKUP := os.Getenv("DB_DSN_BACKUP")
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(DB_DSN_BACKUP),
	))
	DB = bun.NewDB(sqlDB, pgdialect.New())

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	fmt.Println("Connected to the database2.")
	// クエリーフックを追加
	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
		return
	}

	fmt.Println("Migration script executed.")
}
