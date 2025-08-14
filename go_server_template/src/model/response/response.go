package response

import (
	"template/utils"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code        utils.ErrorType `json:"code"`
	Message     string          `json:"message"`
	Success     bool            `json:"success"`
	ResultValue interface{}     `json:"resultValue,omitempty"`
}

type ResponsePage struct {
	*Response
	Total int `json:"total"`
}

type Status struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func ResultPage(data interface{}, total int, c *gin.Context) {
	page := &ResponsePage{
		Total: total,
		Response: &Response{
			utils.Normal,
			utils.CodeDef[utils.Normal],
			true,
			data,
		},
	}
	c.JSON(http.StatusOK, page)
}

func Result(code utils.ErrorType, message string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		message,
		code == utils.Normal,
		data,
	})
}

func ResultString(code utils.ErrorType, message string, data interface{}) ([]byte, error) {
	result := &Response{
		code,
		message,
		code == utils.Normal,
		data,
	}
	return json.Marshal(result)
}

func Ok(c *gin.Context) {
	Result(utils.Normal, utils.CodeDef[utils.Normal], map[string]interface{}{}, c)
}

func Data(data []byte, c *gin.Context) {
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func DataWithContentType(data []byte, contentType string, c *gin.Context) {
	c.Data(http.StatusOK, contentType, data)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(utils.Normal, message, map[string]interface{}{}, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(utils.Normal, utils.CodeDef[utils.Normal], data, c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(utils.Normal, message, data, c)
}

func Fail(code utils.ErrorType, c *gin.Context) {
	Result(code, utils.CodeDef[code], map[string]interface{}{}, c)
}

func FailWithMessage(code utils.ErrorType, message string, c *gin.Context) {
	Result(code, message, map[string]interface{}{}, c)
}

func FailWithDetailed(code utils.ErrorType, data interface{}, c *gin.Context) {
	Result(code, utils.CodeDef[code], data, c)
}
