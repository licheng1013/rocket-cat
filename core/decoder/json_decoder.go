package decoder

import (
	"core/message"
)

type JsonDecoder struct {
}

// EncodeBytes 编码为字节
func (d JsonDecoder) EncodeBytes(result interface{}) []byte {
	return result.([]byte)
}

func NewJsonDecoder() *JsonDecoder {
	return &JsonDecoder{}
}

// DecoderBytes 处理客户端返回的数据
func (d JsonDecoder) DecoderBytes(bytes []byte) message.Message {
	json := message.JsonMessage{}
	// 这里转换成了map
	message.MsgKit.BytesToStruct(bytes, &json)
	return &json
}
