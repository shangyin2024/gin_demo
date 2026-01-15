package response

import (
	"gin_demo/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    Code   `json:"code"`            // 业务状态码
	Message string `json:"message"`         // 提示信息
	Data    any    `json:"data,omitempty"`  // 响应数据
	Error   string `json:"error,omitempty"` // 错误详情
}

// Success 成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeOK,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	// 提取业务错误
	var bizErr *errors.Error
	if e, ok := err.(*errors.Error); ok {
		bizErr = e
	} else {
		bizErr = Wrap(err, CodeInternalError, "服务器内部错误")
	}

	// HTTP 状态码映射
	httpStatus := getHTTPStatus(bizErr.Code)

	// 构建响应
	resp := Response{
		Code:    bizErr.Code,
		Message: bizErr.Message,
	}

	// 开发环境返回详细错误
	if gin.Mode() == gin.DebugMode && bizErr.Err != nil {
		resp.Error = bizErr.Err.Error()
	}

	c.JSON(httpStatus, resp)
}

// ErrorWithCode 指定错误码的错误响应
func ErrorWithCode(c *gin.Context, code Code, message string) {
	c.JSON(getHTTPStatus(code), Response{
		Code:    code,
		Message: message,
	})
}

// getHTTPStatus 根据业务错误码映射 HTTP 状态码
func getHTTPStatus(code Code) int {
	switch code {
	case CodeOK:
		return http.StatusOK
	case CodeInvalidParams:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeAlreadyExists:
		return http.StatusConflict
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	case CodeInvalidPassword:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
