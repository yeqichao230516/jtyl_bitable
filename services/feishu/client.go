package services_feishu

import "net/http"

type Client struct {
	appID             string
	appSecret         string
	httpClient        *http.Client
	tenantAccessToken string
}
