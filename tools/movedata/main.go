package main

import (
	"fmt"

	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"github.com/simesaba80/toybox-back/pkg/db"
	"github.com/simesaba80/toybox-back/tools/movedata/internal/connect"
)

func main() {
	// This is a placeholder main function.
	// The actual migration logic would be implemented here.
	config.LoadEnv()
	fmt.Println("Environment variables loaded.")
	//移行後のDB接続
	db.Init()
	fmt.Println("Connected to the database1.")
	//移行前のDB接続
	connect.Connect()
	fmt.Println("Connected to the database2.")

	fmt.Println("Migration script executed.")
}
