package message

type Message interface {
	GetMerge() int64
	GetBody() []byte
	GetHeartbeat() bool
	GetCode() int64
	GetMessage() string
	// GetBytesResult 返回字节数据
	GetBytesResult() []byte
	// SetBody 设置消息
	SetBody([]byte) Message
	// GetHeaders 用于扩展其他参数
	GetHeaders() string
	// Bind 绑定到对象上
	Bind(v interface{}) (err error)
}
