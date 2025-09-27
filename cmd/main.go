package main

import (
	"cruder/internal/config"
	"cruder/internal/controller"
	"cruder/internal/handler"
	"cruder/internal/middleware"
	"cruder/internal/repository"
	"cruder/internal/service"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Failed to load .env file")
		os.Exit(1)
	}

	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		logger.Error("X_API_KEY environment variable is not set")
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		logger.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	dsn, err := cfg.GetDSN()
	if err != nil {
		logger.Error("failed to get database connection string", slog.Any("err", err))
		os.Exit(1)
	}

	dbConn, err := repository.NewPostgresConnection(dsn)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("err", err))
		os.Exit(1)
	}

	repositories := repository.NewRepository(dbConn.DB())
	services := service.NewService(repositories)
	controllers := controller.NewController(services)

	loggerMiddleware := middleware.NewLoggerMiddleWare(logger)
	apiKeyMiddleware := middleware.NewApiKeyMiddleware(apiKey)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(loggerMiddleware.Handler())
	r.Use(apiKeyMiddleware.Handler())
	r.SetTrustedProxies(nil)

	handler.New(r, controllers.Users)
	if err := r.Run(); err != nil {
		logger.Error("failed to run server", slog.Any("err", err))
	}
}

func loadConfig() (*config.Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return &cfg, nil
}
