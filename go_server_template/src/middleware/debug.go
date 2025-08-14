package middleware

import (
	"template/global"
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func Debug() gin.HandlerFunc {
	return func(c *gin.Context) {
		info := fmt.Sprintf("request query : %s", c.Request.URL.Path)
		if len(c.Request.URL.RawQuery) > 0 {
			info = fmt.Sprintf("%s?%s", info, c.Request.URL.RawQuery)
		}
		data, _ := c.GetRawData()
		if len(data) > 0 {
			info = fmt.Sprintf("%s, data : %s", info, string(data))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
		global.SYS_LOG.Debug(info)
		c.Next()
	}
}
