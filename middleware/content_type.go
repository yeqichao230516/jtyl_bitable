package middleware

import "github.com/gin-gonic/gin"

// ContentTypeJSON 中间件，用于为响应设置JSON内容类型
func ContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先处理请求（执行业务逻辑）
		c.Next()

		// 业务逻辑执行完成后，检查并设置响应头
		// 1. 获取当前已设置的Content-Type
		contentType := c.Writer.Header().Get("Content-Type")

		// 2. 如果未设置Content-Type，则设置为application/json
		if contentType == "" {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
		// 也可以选择强制覆盖所有响应，但通常不建议：
		// c.Header("Content-Type", "application/json; charset=utf-8")
	}
}
