package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `yaml:"environment" env:"APP_ENV" env-default:"development"`
	GRPC        struct {
		Address         string        `yaml:"address" env:"GRPC_ADDRESS" env-default:":50051"`
		MaxReceiveBytes int64         `yaml:"max_receive_bytes" env:"GRPC_MAX_RECEIVE_BYTES" env-default:"4194304"`
		MaxSendBytes    int64         `yaml:"max_send_bytes" env:"GRPC_MAX_SEND_BYTES" env-default:"4194304"`
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"GRPC_SHUTDOWN_TIMEOUT" env-default:"10"`
		TLSCertFile     string        `yaml:"tls_cert_file" env:"GRPC_TLS_CERT_FILE"`
		TLSKeyFile      string        `yaml:"tls_key_file" env:"GRPC_TLS_KEY_FILE"`

		KeepAliveTime       time.Duration `yaml:"keepalive_time" env:"GRPC_KEEPALIVE_TIME" env-default:"2h"`
		KeepAliveTimeout    time.Duration `yaml:"keepalive_timeout" env:"GRPC_KEEPALIVE_TIMEOUT" env-default:"20s"`
		EnforcementMinTime  time.Duration `yaml:"enforcement_min_time" env:"GRPC_ENFORCEMENT_MIN_TIME" env-default:"5s"`
		PermitWithoutStream bool          `yaml:"permit_without_stream" env:"GRPC_PERMIT_WITHOUT_STREAM" env-default:"false"`
	}
	HTTP struct {
		Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:":9090"`
	}
	Mongo struct {
		URL        string        `yaml:"url" env:"MONGO_URL" env-default:"mongodb://localhost:27017"`
		Database   string        `yaml:"database" env:"MONGO_DATABASE" env-default:"products"`
		Collection string        `yaml:"collection" env:"MONGO_PRODUCT_COLLECTION" env-default:"products"`
		Timeout    time.Duration `yaml:"timeout" env:"MONGO_TIMEOUT" env-default:"10"`
	}
}

func Load(path string, logger *slog.Logger) (Config, error) {
	// env.ReadConfig reads environment variables into the Config struct
	// env.With
	var cfg Config
	var err error
	if path == "" {
		err = cleanenv.ReadEnv(&cfg)
		logger.Info("config loaded from environment")
	} else {
		err = cleanenv.ReadConfig(path, &cfg)
		logger.Info("config loaded from file", "path", path)
	}
	if err != nil {
		return Config{}, fmt.Errorf("load config error: %w", err)
	}
	return cfg, nil
}
