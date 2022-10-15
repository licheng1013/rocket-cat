package message

var message interface{} = DefaultMessage{}

type Message interface {
	GetMerge() int64
}

// DefaultMessage 必须实现 Message 接口
// Json处理则必须先转换为json才能继续处理其他东西
type DefaultMessage struct {
	Merge int64
	Body  []byte
}

func (d DefaultMessage) GetMerge() int64 {
	return d.Merge
}

// GetMessage  获取消息解析
func GetMessage() interface{} {
	return message
}
