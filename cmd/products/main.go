package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/adapters/repository/mongodb"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("product service is stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	configPath := flag.String("config", "", "optional YAML configuration file")
	flag.Parse()
	cfg, err := config.Load(*configPath, logger)
	if err != nil {
		return err
	}

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	connectCtx, cancel := context.WithTimeout(rootCtx, cfg.Mongo.Timeout)
	defer cancel()
	mongoRepo, err := mongodb.New(connectCtx, cfg.Mongo.URL, cfg.Mongo.Database, cfg.Mongo.Collection, logger)
	if err != nil {
		return err
	}

	logger.Info("mongo object", "mongo_obj", mongoRepo)
	return nil
}
