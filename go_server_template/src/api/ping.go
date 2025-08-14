package api

import (
	"github.com/gin-gonic/gin"
)

// @Tags /template
// @Summary  template ping
// @Description template ping
// @Security ApiKeyAuth
// @Param Authorization header string true "token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJcCI6IjI3LjE5Ni4yMDMuODkiLCJleHAiOjE3MTkyMTcyNTQsImlhdCI6MTcxOTIxMDA1NCwiaXNzIjoid3d3LmFlZ2lzLmNvbSIsImp0aSI6IkRHRUhPS1BCQ1IiLCJzZXNzaW9uSWQiOiI1ZDAwOThkMDFmZDU0NGVmOWY3YWJjZDQwMTAwMDBmYiIsInN1YiI6IjAyZDNlOGMyOWYzYjQxMmRhMzIzNjkyZTYxYmVmNGIzIiwic3ViVHlwZSI6InVzZXIiLCJ0b2tlblRUTCI6NzIwMDAwMCwidXNlck5hbWUiOiJzZ2l0Z196aHVqdW55aSJ9.Z8_2yZDhi0jfk7R51z4ZF1GAFMssj657AyhtqH87jf4)
// @Router /template/ping [get]
func Ping(c *gin.Context) {
	c.Writer.WriteString("pong")
}
