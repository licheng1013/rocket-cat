package connect

type Socket interface {
	// ListenBack 监听连接收到的消息，回写到上层方法，当返回 byte 不为空时则写入到客户端
	ListenBack(func([]byte) []byte)
	ListenAddr(addr string)
}

type FishType int

// Addr ----------------- 这里时测试数据
const Addr = "127.0.0.1:12345"
const HelloMsg = "HelloWorld"
