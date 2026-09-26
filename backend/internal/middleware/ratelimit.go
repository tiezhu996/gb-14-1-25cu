package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/util"
)

// RateLimit 限流中间件：基于 Redis 滑动窗口（IP + 路径）。
func RateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil {
			c.Next()
			return
		}
		key := "rl:" + c.ClientIP() + ":" + c.Request.URL.Path
		ctx := context.Background()
		n, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if n == 1 {
			rdb.Expire(ctx, key, window)
		}
		if n > int64(limit) {
			util.Fail(c, http.StatusTooManyRequests, constants.CodeRateLimited, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
