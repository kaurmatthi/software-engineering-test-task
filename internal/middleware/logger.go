package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggerMiddleware struct {
	logger *slog.Logger
}

func NewLoggerMiddleWare(logger *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

func (lm *LoggerMiddleware) Handler() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()

		attrs := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.Duration("latency_ms", time.Duration(time.Since(start).Milliseconds())),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("status", status),
		}

		args := make([]any, len(attrs))
		for i, a := range attrs {
			args[i] = a
		}

		switch {
		case status >= 500:
			lm.logger.Error("http_request", args...)
		case status >= 400:
			lm.logger.Warn("http_request", args...)
		default:
			lm.logger.Info("http_request", args...)
		}
	}
}
