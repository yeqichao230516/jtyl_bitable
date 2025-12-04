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
	r.GET("/ping", api.Ping)
	protected := r.Group("/jtcw_bitable")
	protected.Use(middleware.BearerToken())
	{
		protected.POST("/customer_details", api.PostCustomerDetails)

	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", global.CONFIG.Addr.Host, global.CONFIG.Addr.Port),
		Handler: r,
	}
}
