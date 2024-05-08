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
	OnClose(close func(socketId uint32))
}

type SocketClose interface {
	OnClose(socketId uint32)
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
	ProxyMethod func(socketId uint32, message []byte) []byte //代理方法
	UuidOnCoon  sync.Map                                     // 连接
	//queue       chan []byte                              //结果
	Pool    *common.Pool          //线程池
	onClose func(socketId uint32) //关闭钩子，当链接关闭时触发
	Tls     *Tls
	Debug   bool // 是否开启调试模式，开启后会打印发送日志
}

type Tls struct {
	CertFile string //证书文件
	KeyFile  string //密钥文件
}

// InvokeMethod 此处添加至线程池进行远程调用
func (s *MySocket) InvokeMethod(socketId uint32, message []byte) {
	s.Pool.AddTask(func() {
		method := s.ProxyMethod(socketId, message)
		value, ok := s.UuidOnCoon.Load(socketId)
		if ok {
			s.sendChan(value, method)
			//value.(chan []byte) <- socket.proxyMethod(socketId, message)
		}
	})
}

// AsyncResult 这里是同步的，因为流不允许并发写入
func (s *MySocket) AsyncResult(socketId uint32, f func(bytes []byte)) {
	go func() {
		for {
			value, ok := s.UuidOnCoon.Load(socketId)
			if ok {
				// 使用select语句判断chan是否已经关闭
				select {
				case bytes, v := <-value.(chan []byte):
					if v {
						if s.Debug {
							common.RocketLog.Println("发送数据:", string(bytes))
						}
						if bytes == nil || len(bytes) == 0 {
							continue // 返回的数据为空则不写入客户端
						}
						f(bytes)
					} else {
						return
					}
				}
			} else { // 修复循环问题
				return
			}
		}
	}()
}

// SendSelectMessage 选择id发送
func (s *MySocket) SendSelectMessage(bytes []byte, socketIds ...uint32) {
	for _, item := range socketIds {
		value, ok := s.UuidOnCoon.Load(item)
		if ok {
			value.(chan []byte) <- bytes
		}
	}
}

// SendMessage 广播功能
func (s *MySocket) SendMessage(bytes []byte) {
	s.UuidOnCoon.Range(func(key, value any) bool {
		//value.(chan []byte) <- bytes
		s.sendChan(value, bytes)
		return true
	})
}

// 判断 chan 是否关闭,没关闭则发送
func (s *MySocket) sendChan(value any, data []byte) {
	if value == nil {
		return
	}
	value.(chan []byte) <- data
}

// OnClose 不允许使用此钩子注册,因为会被覆盖,请使用插件并实现SocketClose接口
func (s *MySocket) OnClose(close func(socketId uint32)) {
	s.onClose = close
}

// close 关闭连接并处理通道
func (s *MySocket) close(socketId uint32) {
	_, ok := s.UuidOnCoon.Load(socketId)
	if ok {
		s.UuidOnCoon.Delete(socketId)
		//close(value.(chan []byte))
	}
	if s.onClose != nil {
		s.onClose(socketId)
	}
}

//func (s *MySocket) SendGatewayMessage(bytes []byte) {
//	panic("字类没有重新实现广播!")
//}

// 获取一个新chan和uuid
func (s *MySocket) getNewChan() (socketId uint32) {
	for {
		socketId = common.UuidKit.UUID()
		_, ok := s.UuidOnCoon.Load(socketId)
		if !ok { // 如果不存在则创建
			s.UuidOnCoon.Store(socketId, make(chan []byte, 120)) // 通道缓冲区大小为120)
			return
		}
	}
}

// 处理err,如果err不为空则关闭连接
func (s *MySocket) handleErr(err error, socketId uint32, errInfo string) bool {
	if err != nil {
		common.FileLogger().Println(errInfo + err.Error())
		s.close(socketId)
		return true
	}
	return false
}
