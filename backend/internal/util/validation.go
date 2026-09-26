package util

import "github.com/blueship581/codelearn/internal/constants"

// ValidationCode 返回参数校验错误码（handler 统一使用）。
func ValidationCode() int {
	return constants.CodeValidation
}

// ValidationMessage 返回参数校验错误文案。
func ValidationMessage() string {
	return constants.MsgValidationError
}
