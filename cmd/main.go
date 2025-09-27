package main

import (
	"cruder/internal/controller"
	"cruder/internal/handler"
	"cruder/internal/middleware"
	"cruder/internal/repository"
	"cruder/internal/service"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	dsn := os.Getenv("POSTGRES_DSN")
	apiKey := os.Getenv("X_API_KEY")
	if apiKey == "" {
		logger.Error("X_API_KEY environment variable is not set")
		os.Exit(1)
	}
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
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

	handler.New(r, controllers.Users)
	if err := r.Run(); err != nil {
		logger.Error("failed to run server", slog.Any("err", err))
	}
}
