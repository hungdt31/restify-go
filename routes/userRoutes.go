package routes

import (
	"myapp/controllers"
	"myapp/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := NewBaseRoute(r, "/users").Group()
	{
		userGroup.POST("", controllers.CreateUser)
		userGroup.GET("", middleware.AuthRequired(), controllers.GetUsers)
	}
}
