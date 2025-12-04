package global

import (
	"jtyl_bitable/model"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	RECORDS_ID []string
	LOGGER     *logrus.Logger
	VIPER      *viper.Viper
	CONFIG     *model.Config
	FEISHU     *lark.Client
	HTTP       *http.Server
)
