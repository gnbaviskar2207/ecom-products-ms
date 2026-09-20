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
	"github.com/gnbaviskar2207/ecom-products-ms/internal/services"
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
	configPath := flag.String("config", "./configs/dev/config.yaml", "optional YAML configuration file")
	flag.Parse()
	logger.Info("config path", "configPath", *configPath)
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
	productService := services.ProductServiceNew(mongoRepo)
	logger.Info("product service object", "service", productService)

	// Add shutdown logic
	go func() {
		<-rootCtx.Done()
		logger.Info("shutting down product service")
		mongoRepo.Close(rootCtx)
		cancel()
	}()

	return nil
}
