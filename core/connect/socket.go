package connect

import "sync"

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
	proxyMethod func([]byte) []byte
	uuidOnCoon  sync.Map
}
