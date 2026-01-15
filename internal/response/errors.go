package response

import "gin_demo/pkg/errors"

// Code 业务错误码（项目特定）
type Code = errors.Code

const (
	// 成功
	CodeOK Code = 0

	// 客户端错误 (10xxx)
	CodeInvalidParams   Code = 10001 // 参数错误
	CodeUnauthorized    Code = 10002 // 未授权
	CodeForbidden       Code = 10003 // 禁止访问
	CodeNotFound        Code = 10004 // 资源不存在
	CodeAlreadyExists   Code = 10005 // 资源已存在
	CodeTooManyRequests Code = 10006 // 请求过于频繁
	CodeInvalidPassword Code = 10007 // 密码错误

	// 服务端错误 (50xxx)
	CodeInternalError Code = 50001 // 内部错误
	CodeDatabaseError Code = 50002 // 数据库错误
	CodeCacheError    Code = 50003 // 缓存错误
)

// 错误码消息映射
var codeMessages = map[Code]string{
	CodeOK:              "success",
	CodeInvalidParams:   "参数错误",
	CodeUnauthorized:    "未授权",
	CodeForbidden:       "禁止访问",
	CodeNotFound:        "资源不存在",
	CodeAlreadyExists:   "资源已存在",
	CodeTooManyRequests: "请求过于频繁",
	CodeInvalidPassword: "密码错误",
	CodeInternalError:   "服务器内部错误",
	CodeDatabaseError:   "数据库错误",
	CodeCacheError:      "缓存错误",
}

// Message 获取错误码对应的消息
func Message(code Code) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// New 创建业务错误
func New(code Code, message string) *errors.Error {
	return errors.New(code, message)
}

// NewWithError 创建带原始错误的业务错误
func NewWithError(code Code, message string, err error) *errors.Error {
	return errors.NewWithError(code, message, err)
}

// Wrap 包装错误
func Wrap(err error, code Code, message string) *errors.Error {
	return errors.Wrap(err, code, message)
}

// 预定义错误（常用错误的快捷方式）
var (
	ErrInvalidParams   = errors.New(CodeInvalidParams, Message(CodeInvalidParams))
	ErrUnauthorized    = errors.New(CodeUnauthorized, Message(CodeUnauthorized))
	ErrForbidden       = errors.New(CodeForbidden, Message(CodeForbidden))
	ErrNotFound        = errors.New(CodeNotFound, Message(CodeNotFound))
	ErrAlreadyExists   = errors.New(CodeAlreadyExists, Message(CodeAlreadyExists))
	ErrTooManyRequests = errors.New(CodeTooManyRequests, Message(CodeTooManyRequests))
	ErrInvalidPassword = errors.New(CodeInvalidPassword, Message(CodeInvalidPassword))
	ErrInternalError   = errors.New(CodeInternalError, Message(CodeInternalError))
	ErrDatabaseError   = errors.New(CodeDatabaseError, Message(CodeDatabaseError))
	ErrCacheError      = errors.New(CodeCacheError, Message(CodeCacheError))
)
