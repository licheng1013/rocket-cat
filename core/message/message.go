package message

var message interface{} = JsonMessage{}

type Message interface {
	GetMerge() int64
}

// GetMessage  获取消息解析
func GetMessage() interface{} {
	return message
}
