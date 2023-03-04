package common

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// 断言工具

func AssertErr(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertPtr(v interface{}, errInfo string) (err error) {
	pv1 := reflect.ValueOf(v)
	if pv1.Kind() != reflect.Ptr {
		err = errors.New(errInfo)
		return
	}
	return nil
}

func AssertNil(v interface{}, errorInfo string) {
	if v == nil {
		panic(errorInfo)
	}
}

const (
	GatewayName  = "gateway-game"
	ServicerName = "service-game"
)

// Closeing 关闭接口
type Closeing interface {
	Close()
}

func StopApplication() {
	log.Println("监听关机中...")
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                                             // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("正在关机中...")
}
