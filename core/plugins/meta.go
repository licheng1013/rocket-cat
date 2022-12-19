package plugins

import (
	"core/message"
	"net"
	"sync"
)

// Meta 链接信息
type Meta struct {
	// 链接的Id
	SessionId uint32
	// 链接
	Conn net.Conn
	// 传入消息
	Message message.Message
	// 所有链接
	App        []net.Conn
	// 心跳超时管理
	TimeOutMap *sync.Map
}
