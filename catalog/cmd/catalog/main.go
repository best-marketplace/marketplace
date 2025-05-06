package main

import (
	"catalog/internal/app"
	"catalog/internal/config"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	log := setupLogger()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application := app.NewApp(log, cfg)
	go func() {
		if err := application.Run(); err != nil {
			log.Error(fmt.Errorf("error running server %w", err).Error())
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error(fmt.Errorf("error during shutdown: %w", err).Error())
	}

	log.Info("Server stopped")

}

// setupLogger инициализирует логгер
func setupLogger() *slog.Logger {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
