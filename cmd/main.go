package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"subscription_aggregator_go/internal/app"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	defer a.Logger.Sync()

	go func() {
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	a.Logger.Info("server exited gracefully")
}
