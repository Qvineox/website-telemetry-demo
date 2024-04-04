package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"path/filepath"
)

type StaticConfig struct {
	Host string `env-required:"true" env:"host" json:"host"`
	Port string `env-required:"true" env:"port" json:"port"`

	AllowedOrigins string `env:"allowed_origins" json:"allowed_origins"`
}

func NewStaticConfig() (StaticConfig, error) {
	slog.Info("loading static configuration...")

	var cfg StaticConfig

	currentDir, err := os.Getwd()
	if err != nil {
		return StaticConfig{}, err
	}

	// if config file not find tries to get configuration parameters from environment
	err = cleanenv.ReadConfig(filepath.Join(currentDir, "configs", "config.static.json"), &cfg)
	if err != nil {
		slog.Warn("config file not found, reading environment...")
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			return StaticConfig{}, err
		}
	}

	slog.Info("static configuration loaded")
	return cfg, nil
}
