package v1

import (
	"github.com/gin-gonic/gin"
	"template/dao"
	"template/global"
	"template/model/request"
	"template/model/response"
	"template/utils"
)

func ImportCount(c *gin.Context) {
	param := request.ImportCount{}
	if err := c.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage(utils.InvaildParams, "parse param error", c)
		return
	}
	count, err := dao.ImportCount(&param)
	if err != nil {
		global.SYS_LOG.Error(err.Error())
		response.FailWithMessage(utils.GetType(err), err.Error(), c)
		return
	}
	response.OkWithData(count, c)
}
