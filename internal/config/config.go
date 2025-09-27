package config

import (
	"errors"
	"os"
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
