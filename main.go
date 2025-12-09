package main

import (
	"jtyl_bitable/core"
	"jtyl_bitable/global"
	"jtyl_bitable/initialize"
)

func main() {
	initializeSystem()
	core.RunServer()
}
func initializeSystem() {
	global.LOGGER = initialize.Logger()
	global.VIPER = initialize.Viper()
	global.CONFIG = initialize.Config()
	global.FEISHU = initialize.FeiShu()
	global.HOOK = initialize.Hook()
	global.HTTP = initialize.Http()
}
