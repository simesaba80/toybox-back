package db

import (
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/simesaba80/toybox-back/pkg/config"
)

var DB *bun.DB

func Init() {
	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(config.DB_DSN),
	))
	DB = bun.NewDB(sqlDB, pgdialect.New())

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}

	// クエリーフックを追加
	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))
}
