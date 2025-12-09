package initialize

import (
	"jtyl_bitable/api"
	"jtyl_bitable/global"

	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func Hook() *dispatcher.EventDispatcher {
	return dispatcher.NewEventDispatcher(global.CONFIG.Event.VerificationToken, global.CONFIG.Event.EncryptKey).
		OnCustomizedEvent("approval_instance", api.Handler)
}
