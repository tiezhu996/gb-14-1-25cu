package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/util"
)

// ErrorHandler 全局错误处理中间件：捕获 panic 与业务错误，统一 JSON 响应。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "error", r, "path", c.Request.URL.Path, "request_id", GetRequestID(c))
				util.Fail(c, http.StatusInternalServerError, constants.CodeInternal, constants.MsgInternalError)
				c.Abort()
			}
		}()
		c.Next()
		// 统一处理由 handler 写入的错误响应
		if len(c.Errors) > 0 {
			first := c.Errors[0].Err
			var ae *util.AppError
			if errors.As(first, &ae) {
				status := http.StatusBadRequest
				switch {
				case ae.Code >= 2000 && ae.Code < 3000:
					status = http.StatusBadRequest
				case ae.Code == constants.CodeUnauthorized || ae.Code == constants.CodeInvalidToken || ae.Code == constants.CodeTokenExpired:
					status = http.StatusUnauthorized
				case ae.Code == constants.CodeForbidden:
					status = http.StatusForbidden
				case ae.Code == constants.CodeNotFound || ae.Code == constants.CodeUserNotFound ||
					ae.Code == constants.CodeCourseNotFound || ae.Code == constants.CodeProblemNotFound ||
					ae.Code == constants.CodeSubmissionNotFound || ae.Code == constants.CodeDiscussionNotFound ||
					ae.Code == constants.CodeAchievementNotFound || ae.Code == constants.CodeAuditNotFound:
					status = http.StatusNotFound
				case ae.Code == constants.CodeConflict || ae.Code == constants.CodeUserExists ||
					ae.Code == constants.CodeCourseExists || ae.Code == constants.CodeProblemExists:
					status = http.StatusConflict
				case ae.Code == constants.CodeRateLimited:
					status = http.StatusTooManyRequests
				}
				util.Fail(c, status, ae.Code, ae.Message)
				return
			}
			var verrs validator.ValidationErrors
			if errors.As(first, &verrs) {
				util.Fail(c, http.StatusBadRequest, constants.CodeValidation, constants.MsgValidationError)
				return
			}
			util.Fail(c, http.StatusInternalServerError, constants.CodeInternal, constants.MsgInternalError)
		}
	}
}
