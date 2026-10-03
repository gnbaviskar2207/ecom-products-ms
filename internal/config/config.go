package config

import (
	"log/slog"

	commonConfig "github.com/gnbaviskar2207/ecom-common/pkg/config"
)

type Config struct {
	commonConfig.CommonConfig `yaml:",inline"`
	commonConfig.MongoConfig  `yaml:"mongo"`
	commonConfig.GRPCConfig   `yaml:"grpc"`
	commonConfig.HTTPConfig   `yaml:"http"`
}

func Load(configPath string, logger *slog.Logger) (*Config, error) {
	return commonConfig.Load[Config](configPath, logger)
}
