package config

import (
	"log/slog"

	commonConfig "github.com/gnbaviskar2207/ecom-common/pkg/config"
)

type Config struct {
	commonConfig.CommonConfig `yaml:",inline"`

	// Product ProductConfig `yaml:"product"`
}

func Load(configPath string, logger *slog.Logger) (*Config, error) {
	return commonConfig.Load[Config](configPath, logger)
}
