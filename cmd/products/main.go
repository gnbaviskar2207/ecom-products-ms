package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gnbaviskar2207/ecom-common/pkg/telemetry"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("product service is stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logBaseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(telemetry.NewTracehandler(logBaseHandler))
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
	defer rootCtxCancel()

	shutDownTracer, err := telemetry.InitTracer(rootCtx, telemetry.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.Environment,
		CollectorURL:   "localhost:4317",
	})
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		return err
	}

	defer func() {
		if shutDownTracer != nil {
			err = shutDownTracer(context.Background())
		}
	}()

	connectCtx, connectCancel := context.WithTimeout(rootCtx, srv.cfg.MongoConfig.Timeout)
	defer connectCancel()
	if err := srv.connectMongo(connectCtx); err != nil {
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
