package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	// Graceful shutdown
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	e := app.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer stop()

		log.Println("Server starting on :8080")
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Printf("ERROR: Server stopped unexpectedly: %v", err)
		}
	}()

	// シグナルを待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutdown signal received. Initiating shutdown...")
	case <-ctx.Done():
		log.Println("Server goroutine finished. Initiating shutdown...")
	}

	log.Println("Performing graceful server shutdown...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
			log.Printf("ERROR: Server shutdown failed: %v", err)
		}

	wg.Wait()

	log.Println("Server has shut down. Running final cleanup.")
}