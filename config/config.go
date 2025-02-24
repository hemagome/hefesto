package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Encryption  EncryptionConfig
	Logging     LoggingConfig
}

type EncryptionConfig struct {
	Key string
}

type LoggingConfig struct {
	Level string
	Path  string
}

var cfg *Config

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")

	// Environment variables override
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HEFESTO")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	cfg.Environment = env
	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}
