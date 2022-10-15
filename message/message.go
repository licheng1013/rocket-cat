package message

var message interface{} = DefaultMessage{}

type Message interface {
	GetMerge() int64
}

// DefaultMessage 必须实现 Message 接口
type DefaultMessage struct {
	Merge int64
	Body  string
}

func (d DefaultMessage) GetMerge() int64 {
	return d.Merge
}

// GetMessage  获取消息解析
func GetMessage() interface{} {
	return message
}
