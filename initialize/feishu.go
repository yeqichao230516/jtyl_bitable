package initialize

import (
	"jtyl_bitable/global"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

func FeiShu() *lark.Client {
	return lark.NewClient(global.CONFIG.App.Id, global.CONFIG.App.Secret)
}
