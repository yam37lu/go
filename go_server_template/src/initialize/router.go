package initialize

import (
	"template/global"
	"template/middleware"
	"template/router"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func Routers(contextPath string) *gin.Engine {
	var Router = gin.New()
	Router.Use(middleware.Cors())
	PrivateGroup := Router.Group("")
	{
		router.InitApiRouter(PrivateGroup, contextPath) // 注册功能api路由
	}
	global.SYS_LOG.Info("router register success")
	return Router
}
