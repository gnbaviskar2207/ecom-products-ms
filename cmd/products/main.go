package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	productsV1 "github.com/gnbaviskar2207/ecom-products-ms/gen/products"
	grpcApi "github.com/gnbaviskar2207/ecom-products-ms/internal/adapters/grpc"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/adapters/repository/mongodb"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/config"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/services"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/transform/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
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

	rootCtx, rootCtxCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	connectCtx, connectCancel := context.WithTimeout(rootCtx, cfg.Mongo.Timeout)
	defer connectCancel()
	mongoRepo, err := mongodb.New(connectCtx, cfg.Mongo.URL, cfg.Mongo.Database, cfg.Mongo.Collection, logger)
	if err != nil {
		return err
	}
	transform := &generated.ConverterImpl{}
	productService := services.New(mongoRepo, logger)
	productGRPCAdapter := grpcApi.New(logger, productService, transform)

	listener, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return err
	}

	serverOptions := []grpc.ServerOption{
		// Client ---> [large request] ---> gRPC Server
		//                                  ❌ rejected if > MaxReceiveBytes
		grpc.MaxRecvMsgSize(int(cfg.GRPC.MaxReceiveBytes)),
		//   Server ---> [large response] ---> Client
		//                  ❌ rejected if > MaxSendBytes
		grpc.MaxSendMsgSize(int(cfg.GRPC.MaxSendBytes)),

		// keepalive is used to keep the connection alive
		// Keepalive allows the server to periodically send an HTTP/2 PING frame to check whether the client/connection is still alive.
		grpc.KeepaliveParams(keepalive.ServerParameters{
			// How long the server waits before sending a keepalive PING
			Time: cfg.GRPC.KeepAliveTime,
			// How long the server waits for a response to the PING.
			Timeout: cfg.GRPC.KeepAliveTimeout,
		}),

		// Controls which keepalive PINGs from clients the server will accept.
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			// EnforcementMinTime defines the minimum allowed interval between keepalive PINGs.
			MinTime: cfg.GRPC.EnforcementMinTime,
			// Determines whether the client is allowed to send keepalive PINGs.
			PermitWithoutStream: cfg.GRPC.PermitWithoutStream,
		}),
	}
	serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(
		grpcApi.RecoveryInterceptor(logger),
		grpcApi.ErrorInterceptor(logger),
	))
	grpcServer := grpc.NewServer(serverOptions...)
	productsV1.RegisterProductServiceServer(grpcServer, productGRPCAdapter)
	if cfg.Environment == "development" {
		reflection.Register(grpcServer)
		logger.Info("grpc reflection is enabled")
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	logger.Info("product service is running on", "address", cfg.GRPC.Address)
	select {
	case <-rootCtx.Done():
		logger.Info("shutting down the product server(signal received)")

	case err := <-errCh:
		logger.Error("grpc server stopped unexpectedly", "error", err)
	}

	rootCtxCancel()

	shutDownContext, shutDownCancel := context.WithTimeout(context.Background(), cfg.GRPC.ShutdownTimeout)
	defer shutDownCancel()

	stoppedCh := make(chan struct{})

	go func() {
		// allow the inflight/pending requests to complete
		grpcServer.GracefulStop()
		close(stoppedCh)
	}()

	select {
	case <-shutDownContext.Done():
		logger.Error("graceful shutdown timed out, forcing ")
		grpcServer.Stop()
	case <-stoppedCh:
	}
	logger.Info("grpc server is stopped")
	logger.Info("shutting down mongo db")
	mongodbCloseCtx, mongodbCloseCancel := context.WithTimeout(context.Background(), cfg.Mongo.Timeout)
	defer mongodbCloseCancel()
	if err = mongoRepo.Close(mongodbCloseCtx); err != nil {
		logger.Error("mongodb shutdown failed",
			"error", err,
		)
	} else {
		logger.Info("mongodb gracefull shutdown complete")
	}
	logger.Info("graceful shutdown of product service is complete")
	return nil
}
