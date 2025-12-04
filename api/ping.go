package api

import (
	"jtyl_bitable/model"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, model.Resp{
		Code: 0,
		Msg:  "pong",
	})
}
