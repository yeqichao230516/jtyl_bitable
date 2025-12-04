package initialize

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal("无法打开日志文件:", err)
	}
	logger.SetOutput(file)
	return logger
}
