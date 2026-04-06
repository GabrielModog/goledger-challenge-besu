package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goledger-challenge-besu/config"
	"goledger-challenge-besu/internal/database"
	"goledger-challenge-besu/internal/handler"
	"goledger-challenge-besu/internal/repository"
	"goledger-challenge-besu/internal/router"
	"goledger-challenge-besu/pkg/blockchain"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.NewPostgressConnection(ctx, cfg.ConnectionString)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.RunMigrations(ctx); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	bc, err := blockchain.NewBlockchainClient(
		cfg.RPCURL,
		cfg.ContractAddress,
		cfg.PrivateKey,
		cfg.ABIPath,
	)
	if err != nil {
		slog.Error("Failed to initialize blockchain client", "error", err)
		os.Exit(1)
	}
	defer bc.Close()

	repo := repository.NewStorageRepository(db)
	h := handler.NewStorageHandler(bc, repo)

	app := router.Setup(h)

	go func() {
		slog.Info("Starting server", "port", cfg.ServerPort)
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server stopped")
}
