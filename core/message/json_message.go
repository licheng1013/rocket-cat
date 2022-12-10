package message


// JsonMessage 必须实现 Message 接口
// Json处理则必须先转换为json才能继续处理其他东西
type JsonMessage struct {
	Merge int64
	Body  []byte
}

func (d JsonMessage) GetMerge() int64 {
	return d.Merge
}


