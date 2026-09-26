package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/constants"
)

// Resp 统一响应结构：{ "code": 0, "message": "ok", "data": ... }。
type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Resp{Code: constants.CodeOK, Message: constants.MsgOK, Data: data})
}

// SuccessMessage 返回仅带文案的成功响应。
func SuccessMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Resp{Code: constants.CodeOK, Message: message})
}

// Fail 返回带业务错误码的失败响应。
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Resp{Code: code, Message: message})
}

// PageData 分页数据包装。
type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

// SuccessPage 返回分页成功响应。
func SuccessPage(c *gin.Context, list any, total, page, pageSize int64) {
	Success(c, PageData{List: list, Total: total, Page: page, PageSize: pageSize})
}
