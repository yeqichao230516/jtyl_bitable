package api

import (
	"jtyl_bitable/model"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, model.SuccessResp{
		Code: 200,
		Msg:  "success",
	})
}
