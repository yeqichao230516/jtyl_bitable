package middleware

import "github.com/gin-gonic/gin"

// ContentTypeJSON 中间件，用于为响应设置JSON内容类型
func ContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		contentType := c.Writer.Header().Get("Content-Type")
		if contentType == "" {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
	}
}
