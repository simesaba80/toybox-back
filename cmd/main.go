package main

import (
	"context"
	"log"
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

	// OSからの割り込みシグナルをリッスンするコンテキストを作成します。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		e := app.Start()
		if err := e.Start(":8080"); err != nil {
			log.Printf("ERROR: Server failed to start: %v.", err)
			stop()
		}
	}()

	// シグナルを受信するまでブロックします
	<-ctx.Done()

	log.Println("Process finished.")
}
