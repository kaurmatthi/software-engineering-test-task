package handler

import (
	"cruder/internal/controller"

	_ "cruder/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(router *gin.Engine, userController *controller.UserController) *gin.Engine {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	v1 := router.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			userGroup.GET("/", userController.GetAllUsers)
			userGroup.GET("/username/:username", userController.GetUserByUsername)
			userGroup.GET("/id/:id", userController.GetUserByID)
			userGroup.POST("/", userController.CreateUser)
			userGroup.DELETE("/:id", userController.DeleteUser)
			userGroup.PUT("/:id", userController.UpdateUser)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
