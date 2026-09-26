package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/util"
)

// RequireRole RBAC 权限中间件：仅允许指定角色访问。
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr, _ := role.(string)
		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}
		util.Fail(c, http.StatusForbidden, constants.CodeForbidden, constants.MsgForbidden)
		c.Abort()
	}
}
