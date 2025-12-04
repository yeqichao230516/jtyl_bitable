package initialize

import (
	"jtyl_bitable/global"

	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		global.LOGGER.Fatalf("读取配置文件失败: %v", err)
	}
	return v
}
