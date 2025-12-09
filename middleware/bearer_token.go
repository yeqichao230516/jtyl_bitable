package middleware

import (
	"jtyl_bitable/global"
	"jtyl_bitable/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BearerToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResp{
				Msg: "Authorization header is missing",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResp{
				Msg: "Authorization header format must be Bearer {token}",
			})
			return
		}

		token := parts[1]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResp{
				Msg: "Token not found",
			})
			return
		}

		if token != global.CONFIG.Token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResp{
				Msg: "Invalid token",
			})
			return
		}

		c.Next()
	}
}
