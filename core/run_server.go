package core

import (
	"context"
	"jtyl_bitable/global"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunServer() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := global.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.LOGGER.Errorf("HTTP服务异常: %v", err)
		}
	}()

	<-quit
	global.LOGGER.Info("接收到关闭信号，正在优雅关闭服务...")

	if err := global.HTTP.Shutdown(ctx); err != nil {
		global.LOGGER.Errorf("HTTP服务关闭失败: %v", err)
	}
	wg.Wait()
	global.LOGGER.Info("应用服务关闭，所有资源已释放")
}
