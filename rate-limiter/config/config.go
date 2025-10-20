package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type StorageConfig struct {
	URL  string `mapstructure:"url"`
	Size int    `mapstructure:"size"`
}

type LimiterConfig struct {
	RateLimit     int           `mapstructure:"rate_limit"`
	RateWindow    time.Duration `mapstructure:"rate_window"`
	BlockDuration time.Duration `mapstructure:"block_duration"`
}

type Config struct {
	IP          LimiterConfig            `mapstructure:"ip"`
	Token       map[string]LimiterConfig `mapstructure:"token"`
	StorageType string                   `mapstructure:"storage_type"`
	Storage     map[string]StorageConfig `mapstructure:"storage"`
	ServerPort  string                   `mapstructure:"server_port"`
}

func Load(path, configType string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Enable environment variable support
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
