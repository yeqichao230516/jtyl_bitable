package initialize

import (
	"jtyl_bitable/global"
	"jtyl_bitable/model"
)

func Config() *model.Config {
	return &model.Config{
		Addr: model.ServerAddr{
			Host: global.VIPER.GetString("server_ip"),
			Port: global.VIPER.GetString("server_port"),
		},
		App: model.FeiShuApp{
			Id:     global.VIPER.GetString("app_id"),
			Secret: global.VIPER.GetString("app_secret"),
		},
		Token: global.VIPER.GetString("authorization_bearer_token"),
	}
}
