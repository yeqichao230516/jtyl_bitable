package initialize

import (
	"fmt"
	"jtyl_bitable/api"
	"jtyl_bitable/global"
	"jtyl_bitable/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Http() *http.Server {
	r := gin.Default()
	r.Use(middleware.ContentTypeJSON())
	r.GET("/ping", api.Ping)
	// jtcw := r.Group("/jtcw")
	// jtcw.Use(middleware.BearerToken())
	// {
	// 	jtcw.POST("/customer_details", api.PostCustomerDetails)

	// }
	bltj := r.Group("/bltj")
	bltj.Use(middleware.BearerToken())
	{
		bltj.POST("/performance", api.Performance)

	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", global.CONFIG.Addr.Host, global.CONFIG.Addr.Port),
		Handler: r,
	}
}
