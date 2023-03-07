package decoder

import (
	"github.com/licheng1013/io-game-go/messages"
)

type JsonDecoder struct {
}

// EncodeBytes 编码为字节
func (d JsonDecoder) EncodeBytes(result interface{}) []byte {
	return result.([]byte)
}

// DecoderBytes 处理客户端返回的数据
func (d JsonDecoder) DecoderBytes(bytes []byte) messages.Message {
	json := messages.JsonMessage{}
	// 这里转换成了map
	err := messages.MsgKit.BytesToStruct(bytes, &json)
	if err != nil {
		panic(err)
	}
	return &json
}

// JsonDecoderBytes 工具方法
func JsonDecoderBytes(bytes []byte) messages.Message {
	j := &JsonDecoder{}
	return j.DecoderBytes(bytes)
}
