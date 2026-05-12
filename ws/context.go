package ws

// Context 表示一次消息分发的上下文。
type Context struct {
	Session *Session // 当前会话
	Packet  *Packet  // 当前请求包
}
