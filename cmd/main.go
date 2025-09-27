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
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	dsn := os.Getenv("POSTGRES_DSN")
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

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(loggerMiddleware.Handler())

	handler.New(r, controllers.Users)
	if err := r.Run(); err != nil {
		logger.Error("failed to run server", slog.Any("err", err))
	}
}
