package message

var message interface{} = DefaultMessage{}

type Message interface {
	GetMerge() int
}

// DefaultMessage 必须实现 Message 接口
type DefaultMessage struct {
	Merge int
	Body  string
}

func (d DefaultMessage) GetMerge() int {
	return d.Merge
}

// GetMessage  获取消息解析
func GetMessage() interface{} {
	return message
}
