package initialize

import (
	"fmt"
	"jtyl_bitable/api"
	"jtyl_bitable/global"
	"jtyl_bitable/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
)

func Http() *http.Server {
	r := gin.Default()
	r.Use(middleware.ContentTypeJSON())
	r.GET("/ping", api.Ping)
	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(global.HOOK))

	bltj := r.Group("/bltj")
	bltj.Use(middleware.BearerToken())
	{
		bltj.POST("/performance", api.Performance)
		bltj.POST("/approval", api.CreateApproval)

	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", global.CONFIG.Addr.Host, global.CONFIG.Addr.Port),
		Handler: r,
	}
}
