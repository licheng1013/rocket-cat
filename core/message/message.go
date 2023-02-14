package message

type Message interface {
	GetMerge() int64
	GetBody() []byte
	GetHeartbeat() bool
	GetCode() int64
	GetMessage() string
	GetBytesResult() []byte
	SetBody([]byte)
	// GetHeaders 用于扩展其他参数
	GetHeaders() string
}
