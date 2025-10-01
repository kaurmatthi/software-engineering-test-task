package controller

import "github.com/gin-gonic/gin"

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
