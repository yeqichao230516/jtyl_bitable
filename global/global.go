package global

import (
	"jtyl_bitable/model"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	LOGGER *logrus.Logger
	VIPER  *viper.Viper
	CONFIG *model.Config
	FEISHU *lark.Client
	HTTP   *http.Server
	HOOK   *dispatcher.EventDispatcher
)
