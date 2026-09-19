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
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"GRPC_SHUTDOWN_TIMEOUT" env-default:"10s"`
		TLSCertFile     string        `yaml:"tls_cert_file" env:"GRPC_TLS_CERT_FILE"`
		TLSKeyFile      string        `yaml:"tls_key_file" env:"GRPC_TLS_KEY_FILE"`
	}
	HTTP struct {
		Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:":9090"`
	}
}

func Load(path string) (Config, error) {
	// env.ReadConfig reads environment variables into the Config struct
	// env.With
	var cfg Config
	var err error
	if path == "" {
		err = cleanenv.ReadEnv(&cfg)
		slog.Info("config loaded from environment")
	} else {
		err = cleanenv.ReadConfig(path, &cfg)
		slog.Info("config loaded from file", "path", path)
	}
	if err != nil {
		return Config{}, fmt.Errorf("load config error: %w", err)
	}
	return cfg, nil
}
