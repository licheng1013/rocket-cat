package connect

import (
	"github.com/io-game-go/common"
	"sync"
)

type Socket interface {
	// ListenBack 监听连接收到的消息，回写到上层方法，当返回 byte 不为空时则写入到客户端
	ListenBack(func([]byte) []byte)
	ListenAddr(addr string)
}

// Addr ----------------- 这里时测试数据
const Addr = "192.168.101.10:12345"
const HelloMsg = "HelloWorld"

// MySocket Socket接口的通用字段
type MySocket struct {
	proxyMethod func([]byte) []byte //代理方法
	uuidOnCoon  sync.Map // 连接
	queue       chan []byte //结果
	pool        *common.Pool //线程池
}

// InvokeMethod 此处添加至线程池进行远程调用
func (s *MySocket) InvokeMethod(message []byte) {
	_ = s.pool.AddTaskNonBlocking(func() {
		s.queue <- s.proxyMethod(message)
	})
}
func (s *MySocket) AsyncResult(f func(bytes []byte)) {
	go func() {
		s.queue = make(chan []byte)
		for bytes := range s.queue {
			f(bytes)
		}
	}()
}
