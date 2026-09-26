package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/util"
)

// Auth 认证中间件：解析 Bearer JWT，注入 user_id/role/status 到上下文。
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(secret, tokenStr)
		if err != nil {
			util.Fail(c, http.StatusUnauthorized, constants.CodeInvalidToken, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			util.Fail(c, http.StatusUnauthorized, constants.CodeInvalidToken, constants.MsgUnauthorized)
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Set("role", claims.Role)
		c.Set("status", claims.Status)
		c.Set("username", claims.Username)
		c.Next()
	}
}
