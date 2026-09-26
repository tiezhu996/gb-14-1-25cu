package util

import "fmt"

// AppError 业务错误：携带错误码与可读消息，由 error_handler 中间件统一转换为 JSON 响应。
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

// NewAppError 构造业务错误。
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WrapAppError 包装底层错误并透传业务错误码。
func WrapAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// IsAppError 判断错误链中是否包含指定业务错误码。
func IsAppError(err error, code int) bool {
	for err != nil {
		if ae, ok := err.(*AppError); ok && ae.Code == code {
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}
