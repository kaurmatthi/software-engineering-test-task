package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// helper to create a test router
func setupRouter(apiKey string, ignored []string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(NewApiKeyMiddleware(apiKey, ignored).Handler())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}

// Given: A request without an API key header
func TestApiKeyMiddleware_MissingKey(t *testing.T) {
	r := setupRouter("secret", nil)
	req, _ := http.NewRequest("GET", "/ping", nil)

	// When: The request is sent
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: It should return 401 Unauthorized with a proper error message
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "X-Api-Key header is missing")
}

// Given: A request with an invalid API key header
func TestApiKeyMiddleware_InvalidKey(t *testing.T) {
	r := setupRouter("secret", nil)
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("X-Api-Key", "wrong")

	// When: The request is sent
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: It should return 403 Forbidden with a proper error message
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "provided X-Api-Key is invalid")
}

// Given: A request with a valid API key header
func TestApiKeyMiddleware_ValidKey(t *testing.T) {
	r := setupRouter("secret", nil)
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("X-Api-Key", "secret")

	// When: The request is sent
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: It should return 200 OK and allow the request
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

// Given: A request with an invalid API key header to an ignored endpoint
func TestApiKeyMiddleware_IgnoredEndpoint(t *testing.T) {
	r := setupRouter("secret", []string{"/ping"})
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("X-Api-Key", "wrong")

	// When: The request is sent
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Then: It should return 200 OK and allow the request
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}
