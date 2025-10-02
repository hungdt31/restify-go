package routes

import "github.com/gin-gonic/gin"

// BaseRoute giữ prefix + router engine
type BaseRoute struct {
	Router *gin.Engine
	Prefix string
}

// Hàm tiện ích để đăng ký route có prefix
func (b *BaseRoute) Group() *gin.RouterGroup {
	return b.Router.Group("/api" + b.Prefix)
}

func NewBaseRoute(router *gin.Engine, prefix string) *BaseRoute {
	return &BaseRoute{
		Router: router,
		Prefix: prefix,
	}
}
