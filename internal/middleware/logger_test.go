package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// helper to create a router with logger middleware
func setupLoggerRouter(buf *bytes.Buffer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := slog.NewTextHandler(buf, nil)
	logger := slog.New(handler)
	r := gin.New()
	r.Use(NewLoggerMiddleWare(logger).Handler())

	// endpoints return different status codes
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	r.GET("/bad", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})
	r.GET("/fail", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
	})

	return r
}

func TestLoggerMiddleware_LogsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	// Given: A GET /ok request with query string and User-Agent
	req, _ := http.NewRequest("GET", "/ok?foo=bar", nil)
	req.Header.Set("User-Agent", "test-agent")

	// When: The request is served
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: The log should include key fields
	logOutput := buf.String()

	assert.Contains(t, logOutput, "method=GET")
	assert.Contains(t, logOutput, "path=/ok")
	assert.Contains(t, logOutput, "query=\"foo=bar\"")
	assert.Contains(t, logOutput, "user_agent=test-agent")
	assert.Contains(t, logOutput, "status=200")
	assert.Contains(t, logOutput, "latency_ms=")
}

func TestLoggerMiddleware_InfoLog(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	// Given: A GET /ok request that returns 200 OK
	req, _ := http.NewRequest("GET", "/ok", nil)

	// When: The request is served through the middleware
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: The response is 200 and an INFO log is written
	assert.Equal(t, http.StatusOK, w.Code)
	logOutput := buf.String()
	assert.Contains(t, logOutput, "level=INFO")
	assert.Contains(t, logOutput, "path=/ok")
}

func TestLoggerMiddleware_WarnLog(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	// Given: A GET /bad request that returns 400 Bad Request
	req, _ := http.NewRequest("GET", "/bad", nil)

	// When: The request is served through the middleware
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: The response is 400 and a WARN log is written
	assert.Equal(t, http.StatusBadRequest, w.Code)
	logOutput := buf.String()
	assert.Contains(t, logOutput, "level=WARN")
	assert.Contains(t, logOutput, "path=/bad")
}

func TestLoggerMiddleware_ErrorLog(t *testing.T) {
	var buf bytes.Buffer
	r := setupLoggerRouter(&buf)

	// Given: A GET /fail request that returns 500 Internal Server Error
	req, _ := http.NewRequest("GET", "/fail", nil)

	// When: The request is served through the middleware
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: The response is 500 and an ERROR log is written
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	logOutput := buf.String()
	assert.Contains(t, logOutput, "level=ERROR")
	assert.Contains(t, logOutput, "path=/fail")
}
