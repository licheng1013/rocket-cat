package core

import (
	"core/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Service 新手请不需要之间使用而是 NewService 进行获取对象
type Service struct {
	// 单机模式可用
	router router.Router
}

func (n *Service) SetRouter(router router.Router) {
	n.router = router
}

func NewService() *Service {
	service := &Service{}
	service.SetRouter(&router.DefaultRouter{})
	return service
}

func (n *Service) Stop() {
	log.Println("监听关机中...")
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("正在关机中...")
}
