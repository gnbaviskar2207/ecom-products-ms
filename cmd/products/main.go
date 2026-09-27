package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	configPath := flag.String("config", "./configs/dev/config.yaml", "optional YAML configuration file")
	flag.Parse()
	logger.Info("config path", "configPath", *configPath)
	cfg, err := config.Load(*configPath, logger)
	if err != nil {
		return err
	}
	srv := New(cfg, logger)

	rootCtx, rootCtxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if err := srv.connectMongo(rootCtx); err != nil {
		return err
	}
	metrics := srv.initMetrics()
	srv.buildServices()

	if err := srv.buildGRPCServer(metrics); err != nil {
		return err
	}
	srv.buildHTTPServer()
	srv.registerOpsRoutes()

	srv.start(rootCtx)
	rootCtxCancel()

	shutDownContext, shutDownCancel := context.WithTimeout(context.Background(), srv.cfg.GRPCConfig.ShutdownTimeout)
	defer shutDownCancel()
	return srv.shutDown(shutDownContext)
}
