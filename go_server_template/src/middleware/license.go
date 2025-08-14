package middleware

import (
	"template/global"
	"template/model/response"
	"template/utils"

	"github.com/gin-gonic/gin"
)

func LicenseCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.SYS_LICENSE != nil && !global.SYS_LICENSE.Active {
			response.Fail(utils.ServiceNotAvailable, c)
			c.Abort()
			return
		}
		c.Next()
	}
}
