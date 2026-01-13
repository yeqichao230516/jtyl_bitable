package initialize

import lark "github.com/larksuite/oapi-sdk-go/v3"

func FeiShuClient(appId, appSecret string) *lark.Client {
	return lark.NewClient(appId, appSecret)
}
