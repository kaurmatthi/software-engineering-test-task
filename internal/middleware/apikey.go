package middleware

import "github.com/gin-gonic/gin"

type ApiKeyMiddleware struct {
	apiKey string
}

func NewApiKeyMiddleware(apiKey string) *ApiKeyMiddleware {
	return &ApiKeyMiddleware{apiKey: apiKey}
}

func (am *ApiKeyMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		providedKey := c.GetHeader("X-Api-Key")
		if providedKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "X-Api-Key header is missing"})
			return
		}
		if providedKey != am.apiKey {
			c.AbortWithStatusJSON(403, gin.H{"error": "provided X-Api-Key is invalid"})
			return
		}
		c.Next()
	}
}
