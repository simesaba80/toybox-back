package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/simesaba80/toybox-back/internal/di"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
)

func main() {
	config.LoadEnv()

	// Wireで依存性を注入してアプリケーションを初期化
	app, cleanup, err := di.InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	defer cleanup()

	// Graceful shutdown
	go func() {
		e := app.Start()
		if err := e.Start(":8080"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// シグナルを待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
