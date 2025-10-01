package routes

import (
	"github.com/gin-gonic/gin"
	"myapp/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// middleware chung
	r.Use(middleware.RequestLogger())

	// Load từng group routes
	RegisterUserRoutes(r)

	return r
}
