package connect

import (
	"github.com/licheng1013/rocket-cat/common"
	"sync"
)

type Socket interface {

	// ListenBack 监听连接收到的消息，回写到上层方法，当返回 byte 不为空时则写入到客户端,
	// uuid 为连接建立时的唯一id,message为具体消息内容
	ListenBack(func(uuid uint32, message []byte) []byte)
	ListenAddr(addr string)
	SendSelectMessage(bytes []byte, socketIds ...uint32)
	SendMessage(bytes []byte)
}

//type Broadcast interface {
//	// SendGatewayMessage 发送所有消息
//	SendGatewayMessage(bytes []byte)
//}
//
//type SelectBroadcast interface {
//	// SendSelectMessage 发送指定目标消息
//	SendSelectMessage(bytes []byte, uuid ...uint32)
//}

// Addr ----------------- 这里时测试数据
const Addr = "0.0.0.0:12345"
const HelloMsg = "HelloWorld"

// MySocket Socket接口的通用字段
type MySocket struct {
	proxyMethod func(uuid uint32, message []byte) []byte //代理方法
	UuidOnCoon  sync.Map                                 // 连接
	//queue       chan []byte                              //结果
	Pool    *common.Pool      //线程池
	onClose func(uuid uint32) //关闭钩子，当链接关闭时触发
	Tls     *Tls
	Debug   bool // 是否开启调试模式，开启后会打印发送日志
}

type Tls struct {
	CertFile string //证书文件
	KeyFile  string //密钥文件
}

// InvokeMethod 此处添加至线程池进行远程调用
func (socket *MySocket) InvokeMethod(uuid uint32, message []byte) {
	_ = socket.Pool.AddTaskNonBlocking(func() {
		value, ok := socket.UuidOnCoon.Load(uuid)
		if ok {
			value.(chan []byte) <- socket.proxyMethod(uuid, message)
		}
	})
}

// AsyncResult 这里是同步的，因为流不允许并发写入
func (socket *MySocket) AsyncResult(uuid uint32, f func(bytes []byte)) {
	go func() {
		value, ok := socket.UuidOnCoon.Load(uuid)
		if ok {
			for {
				// 使用select语句判断chan是否已经关闭
				select {
				case bytes, v := <-value.(chan []byte):
					if v {
						if socket.Debug {
							common.Logger().Println("发送数据:", string(bytes))
						}
						if bytes == nil || len(bytes) == 0 {
							continue // 返回的数据为空则不写入客户端
						}
						f(bytes)
					} else {
						return
					}
				}
			}
		}
	}()
}

// SendSelectMessage 选择id发送
func (socket *MySocket) SendSelectMessage(bytes []byte, socketIds ...uint32) {
	for _, item := range socketIds {
		value, ok := socket.UuidOnCoon.Load(item)
		if ok {
			value.(chan []byte) <- bytes
		}
	}
}

// SendMessage 广播功能
func (socket *MySocket) SendMessage(bytes []byte) {
	socket.UuidOnCoon.Range(func(key, value any) bool {
		value.(chan []byte) <- bytes
		return true
	})
}

func (socket *MySocket) OnClose(close func(uuid uint32)) {
	socket.onClose = close
}

// close 关闭连接并处理通道
func (socket *MySocket) close(uuid uint32) {
	value, ok := socket.UuidOnCoon.Load(uuid)
	if ok {
		socket.UuidOnCoon.Delete(uuid)
		close(value.(chan []byte))
	}
	if socket.onClose != nil {
		socket.onClose(uuid)
	}
}

// 初始化线程池
func (socket *MySocket) init() {
	if socket.Pool == nil {
		// 创建一个线程池，指定工作协程数为3，任务队列大小为10
		socket.Pool = common.NewPool(20, 30)
	}
	socket.Pool.Start()
}

//func (s *MySocket) SendGatewayMessage(bytes []byte) {
//	panic("字类没有重新实现广播!")
//}
