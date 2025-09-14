package connect

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func Connect() {
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

}
