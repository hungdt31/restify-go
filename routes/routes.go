package routes

import (
	"github.com/gin-gonic/gin"
	"myapp/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// middleware chung
	r.Use(middleware.RequestLogger())

	// Health check endpoint cho Docker healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Application is running",
		})
	})

	// Load từng group routes
	RegisterUserRoutes(r)
	RegisterAuthRoutes(r)

	return r
}
