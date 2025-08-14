package core

import (
	"template/global"
	"template/initialize"
	"time"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer(address, contextPath string) {
	Router := initialize.Routers(contextPath)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.SYS_LOG.Info("server run success on ", zap.String("address", address))

	global.SYS_LOG.Error(s.ListenAndServe().Error())
}
