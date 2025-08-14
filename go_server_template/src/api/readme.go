package api

import (
	"net/http"
	"template/service"

	"github.com/gin-gonic/gin"
)

func Readme(c *gin.Context) {
	c.JSON(http.StatusOK, service.Readme())
}
