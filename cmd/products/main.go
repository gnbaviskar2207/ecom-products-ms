package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("product service is stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "", "optional YAML configuration file")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}
	slog.Info("config loaded", "cfg", cfg)
	return nil
}
