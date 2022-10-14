package message

var message Message = DefaultMessage{}

type Message interface {
	GetMerge() int
}

type DefaultMessage struct {
	Merge int
	Body  string
}

func (d DefaultMessage) GetMerge() int {
	return d.Merge
}

// GetMessage  获取消息解析
func GetMessage() any {
	return message
}

// SetMessage  设置消息解析
func SetMessage(v Message) {
	message = v
}
