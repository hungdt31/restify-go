package routes

import (
	"github.com/gin-gonic/gin"
	"myapp/controllers"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/", controllers.GetUsers)
	}
}
