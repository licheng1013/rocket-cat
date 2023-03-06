package decoder

import (
	"github.com/licheng1013/io-game-go/message"
)

type JsonDecoder struct {
}

// EncodeBytes 编码为字节
func (d JsonDecoder) EncodeBytes(result interface{}) []byte {
	return result.([]byte)
}

// DecoderBytes 处理客户端返回的数据
func (d JsonDecoder) DecoderBytes(bytes []byte) message.Message {
	json := message.JsonMessage{}
	// 这里转换成了map
	err := message.MsgKit.BytesToStruct(bytes, &json)
	if err != nil {
		panic(err)
	}
	return &json
}

// JsonDecoderBytes 工具方法
func JsonDecoderBytes(bytes []byte) message.Message {
	j := &JsonDecoder{}
	return j.DecoderBytes(bytes)
}
