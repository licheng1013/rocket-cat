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
	uuidOnCoon  sync.Map            // 连接
	queue       chan []byte         //结果
	Pool        *common.Pool        //线程池
}

// InvokeMethod 此处添加至线程池进行远程调用
func (s *MySocket) InvokeMethod(message []byte) {
	_ = s.Pool.AddTaskNonBlocking(func() {
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

// 初始化线程池
func (s *MySocket) init() {
	if s.Pool == nil {
		// 创建一个线程池，指定工作协程数为3，任务队列大小为10
		s.Pool = common.NewPool(20, 30)
	}
	s.Pool.Start()
}
