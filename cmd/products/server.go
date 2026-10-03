package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gnbaviskar2207/ecom-common/pkg/interceptors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	grpcApi "github.com/gnbaviskar2207/ecom-products-ms/internal/adapters/grpc"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/adapters/repository/mongodb"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/config"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/services"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/transform"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/transform/generated"
	prom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	productsV1 "github.com/gnbaviskar2207/ecom-common/pkg/gen/products"
)

type Server struct {
	cfg          *config.Config
	logger       *slog.Logger
	grpcServer   *grpc.Server
	httpServer   *http.Server
	mongoRepo    *mongodb.MongoProductRepository
	listener     net.Listener
	transform    transform.Converter
	service      ports.ProductService
	grpcAdapter  *grpcApi.ProductAdapter
	healthServer *health.Server
	mux          *http.ServeMux
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	return &Server{
		cfg:       cfg,
		logger:    logger,
		transform: &generated.ConverterImpl{},
	}
}

func (s *Server) connectMongo(ctx context.Context) error {
	s.logger.Info("connecting to mongodb", "url", s.cfg.MongoConfig.URL, "database", s.cfg.MongoConfig.Database, "collection", s.cfg.MongoConfig.Collection)
	mongoRepo, err := mongodb.New(ctx, s.cfg.MongoConfig.URL, s.cfg.MongoConfig.Database, s.cfg.MongoConfig.Collection, s.logger)
	if err != nil {
		return err
	}
	s.mongoRepo = mongoRepo
	return nil
}

func (s *Server) buildGRPCServer(metrics *prom.ServerMetrics) error {
	var err error
	s.listener, err = net.Listen("tcp", s.cfg.GRPCConfig.Address)
	if err != nil {
		return err
	}
	serverOptions := []grpc.ServerOption{
		// Otel StatsHandler
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		// Client ---> [large request] ---> gRPC Server
		//                                  ❌ rejected if > MaxReceiveBytes
		grpc.MaxRecvMsgSize(int(s.cfg.GRPCConfig.MaxReceiveBytes)),
		//   Server ---> [large response] ---> Client
		//                  ❌ rejected if > MaxSendBytes
		grpc.MaxSendMsgSize(int(s.cfg.GRPCConfig.MaxSendBytes)),

		// keepalive is used to keep the connection alive
		// Keepalive allows the server to periodically send an HTTP/2 PING frame to check whether the client/connection is still alive.
		grpc.KeepaliveParams(keepalive.ServerParameters{
			// How long the server waits before sending a keepalive PING
			Time: s.cfg.GRPCConfig.KeepAliveTime,
			// How long the server waits for a response to the PING.
			Timeout: s.cfg.GRPCConfig.KeepAliveTimeout,
		}),

		// Controls which keepalive PINGs from clients the server will accept.
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			// EnforcementMinTime defines the minimum allowed interval between keepalive PINGs.
			MinTime: s.cfg.GRPCConfig.EnforcementMinTime,
			// Determines whether the client is allowed to send keepalive PINGs.
			PermitWithoutStream: s.cfg.GRPCConfig.PermitWithoutStream,
		}),
	}

	// TLS configurations
	if s.cfg.TLSCertFile != "" && s.cfg.GRPCConfig.TLSKeyFile != "" {
		tlsCreds, err := credentials.NewServerTLSFromFile(
			s.cfg.GRPCConfig.TLSCertFile,
			s.cfg.GRPCConfig.TLSKeyFile,
		)
		if err != nil {
			return fmt.Errorf("error while loading tls cert and key from files -%w", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCreds))
		s.logger.Info("tls is enabled for grpc server")
	}

	serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(
		metrics.UnaryServerInterceptor(),
		interceptors.RequestLoggerInterceptor(s.logger),
		interceptors.RecoveryInterceptor(s.logger),
		interceptors.ErrorInterceptor(s.logger),
	))
	s.grpcServer = grpc.NewServer(serverOptions...)
	productsV1.RegisterProductServiceServer(s.grpcServer, s.grpcAdapter)

	// health server
	s.healthServer = health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	metrics.InitializeMetrics(s.grpcServer)
	if s.cfg.Environment == "development" {
		reflection.Register(s.grpcServer)
		s.logger.Info("grpc reflection is enabled")
	}
	return nil
}

func (s *Server) buildHTTPServer() *http.Server {
	s.mux = http.NewServeMux()
	httpServer := &http.Server{
		Addr:              s.cfg.HTTPConfig.Address,
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}
	s.httpServer = httpServer
	return s.httpServer
}

func (s *Server) registerOpsRoutes() {
	s.mux.Handle("/metrics", promhttp.Handler())

	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s.mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		pingCtx, pingCancel := context.WithTimeout(r.Context(), s.cfg.MongoConfig.Timeout)
		defer pingCancel()
		if err := s.mongoRepo.Ping(pingCtx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) initMetrics() *prom.ServerMetrics {
	metrics := prom.NewServerMetrics()
	prometheus.MustRegister(metrics)
	return metrics
}

func (s *Server) buildServices() {
	s.service = services.New(s.mongoRepo, s.logger)
	s.grpcAdapter = grpcApi.New(s.logger, s.service, s.transform)
}

func (s *Server) start(ctx context.Context) {
	errCh := make(chan error, 2)
	go func() {
		if err := s.grpcServer.Serve(s.listener); err != nil {
			errCh <- fmt.Errorf("grpc server: %w", err)
		}
	}()
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()

	s.logger.Info("product service is running", "grpc address", s.cfg.GRPCConfig.Address, "http address", s.cfg.HTTPConfig.Address)
	select {
	case <-ctx.Done():
		s.logger.Info("shutting down the product server(signal received)")

	case err := <-errCh:
		s.logger.Error("product server stopped unexpectedly", "error", err)
	}
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

func (s *Server) shutDown(ctx context.Context) error {
	var err error

	_ = s.httpServer.Shutdown(ctx)
	stoppedCh := make(chan struct{})

	go func() {
		// allow the inflight/pending requests to complete
		s.grpcServer.GracefulStop()
		close(stoppedCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Error("graceful shutdown timed out, forcing ")
		s.grpcServer.Stop()
	case <-stoppedCh:
	}

	s.logger.Info("grpc server is stopped")
	s.logger.Info("shutting down mongo db")
	mongodbCloseCtx, mongodbCloseCancel := context.WithTimeout(context.Background(), s.cfg.MongoConfig.Timeout)
	defer mongodbCloseCancel()
	if err = s.mongoRepo.Close(mongodbCloseCtx); err != nil {
		s.logger.Error("mongodb shutdown failed",
			"error", err,
		)
	} else {
		s.logger.Info("mongodb gracefull shutdown complete")
	}
	s.logger.Info("graceful shutdown of product service is complete")
	return err
}
