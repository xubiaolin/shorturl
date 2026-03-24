package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// 硬编码的鉴权 token
	AuthToken = "shorturl-secret-token-2026"
)

// AuthMiddleware 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 获取 token
		token := c.GetHeader("Authorization")

		// 支持 Bearer token 格式
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// 验证 token
		if token == "" || token != AuthToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问，请提供有效的 token",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 设置用户名到上下文
		c.Set("username", "admin")
		c.Next()
	}
}
