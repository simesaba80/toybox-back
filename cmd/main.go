package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/simesaba80/toybox-back/internal/di"
	"github.com/simesaba80/toybox-back/pkg/config"
)

func main() {
	config.LoadEnv()

	// Wireで依存性を注入してアプリケーションを初期化
	app, cleanup, err := di.InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}
	defer cleanup()

	// サーバー起動の失敗を捕捉できるよう、先にシグナルの待機を開始します
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown
	go func() {
		e := app.Start()
		if err := e.Start(":8080"); err != nil {
			log.Printf("ERROR: Server failed to start: %v. Sending SIGTERM to self.", err)
			p, _ := os.FindProcess(os.Getpid())
			p.Signal(syscall.SIGTERM)
		}
	}()

	// シグナルを受信するまでブロックします
	select {
	case sig := <-quit:
		log.Printf("Received signal '%v'. Process will exit.", sig)
	}

	log.Println("Process finished.")
}
