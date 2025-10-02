package routes

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := NewBaseRoute(r, "/auth").Group()
	{
		auth.POST("/login", controllers.Login)
	}
}
