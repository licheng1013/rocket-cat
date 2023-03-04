package connect

import (
	"github.com/io-game-go/common"
	"sync"
)

type Socket interface {

	// ListenBack 监听连接收到的消息，回写到上层方法，当返回 byte 不为空时则写入到客户端,
	// uuid 为连接建立时的唯一id,message为具体消息内容
	ListenBack(func(uuid uint32, message []byte) []byte)
	ListenAddr(addr string)
}

type Broadcast interface {
	SendMessage([]byte)
}

// Addr ----------------- 这里时测试数据
const Addr = "192.168.101.10:12345"
const HelloMsg = "HelloWorld"

// MySocket Socket接口的通用字段
type MySocket struct {
	proxyMethod func(uuid uint32, message []byte) []byte //代理方法
	UuidOnCoon  sync.Map                                 // 连接
	queue       chan []byte                              //结果
	Pool        *common.Pool                             //线程池
}

// InvokeMethod 此处添加至线程池进行远程调用
func (s *MySocket) InvokeMethod(uuid uint32, message []byte) {
	_ = s.Pool.AddTaskNonBlocking(func() {
		s.queue <- s.proxyMethod(uuid, message)
	})
}

// AsyncResult 这里是同步的，因为流不允许并发写入
func (s *MySocket) AsyncResult(f func(bytes []byte)) {
	go func() {
		s.queue = make(chan []byte)
		for bytes := range s.queue {
			if bytes == nil || len(bytes) == 0 {
				continue // 返回的数据为空则不写入客户端
			}
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

//func (s *MySocket) SendMessage(bytes []byte) {
//	panic("字类没有重新实现广播!")
//}
