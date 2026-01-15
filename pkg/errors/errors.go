package errors

import (
	"errors"
	"fmt"
)

// Code 错误码类型（通用）
type Code int

// Error 业务错误（通用结构）
type Error struct {
	Code    Code   // 错误码
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap
func (e *Error) Unwrap() error {
	return e.Err
}

// New 创建新错误
func New(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewWithError 创建带原始错误的业务错误
func NewWithError(code Code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrap 包装错误
func Wrap(err error, code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is 判断是否是指定的错误码
func Is(err error, code Code) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode 获取错误码
func GetCode(err error) Code {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return 0 // 未知错误码
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return err.Error()
}
