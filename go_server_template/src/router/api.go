package router

import (
	"fmt"
	"template/api"
	v1 "template/api/v1"

	// _ "template/docs"
	"template/middleware"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(Router *gin.RouterGroup, contextPath string) {
	Router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{fmt.Sprintf("%s/ping", contextPath)},
		}),
		gin.Recovery(),
	)
	Router.Use(middleware.GinRecovery(true))
	Router.Use(middleware.Cors())
	//Router.Use(middleware.LicenseCheck())
	Router.Use(middleware.Debug())
	commonRouter := Router.Group(contextPath)
	{
		commonRouter.GET("readme", api.Readme)
		commonRouter.GET("ping", api.Ping)
	}
	v1Router := Router.Group(fmt.Sprintf("%s/v1", contextPath))
	//v1Router.Use(middleware.OperatorParser())
	//v1Router.Use(middleware.OperateLog())
	toolRouter := v1Router.Group("/tool")
	{
		toolRouter.POST("count", v1.ImportCount)
	}
}
