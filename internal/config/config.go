package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host    string `mapstructure:"host"`
		DB      string `mapstructure:"db"`
		Port    string `mapstructure:"port"`
		User    string `mapstructure:"user"`
		Sslmode string `mapstructure:"sslmode"`
	}
}

func (c *Config) GetDSN() (string, error) {
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return "", errors.New("POSTGRES_PASSWORD environment variable is not set")
	}
	dsn := "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + password +
		" dbname=" + c.Database.DB +
		" sslmode=" + c.Database.Sslmode

	return dsn, nil
}

func LoadConfig() (*Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return &cfg, nil
}
