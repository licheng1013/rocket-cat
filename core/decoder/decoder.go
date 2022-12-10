package decoder

import (
	"core/message"
	"google.golang.org/protobuf/proto"
)

// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte) (int64, interface{})
}

func ParseResult(result interface{}) []byte {
	// 分发消息
	var bytes []byte
	if result != nil {
		switch result.(type) {
		case []byte:
			bytes = result.([]byte)
			break
		case proto.Message:
			bytes = message.MarshalBytes(result.(proto.Message))
			break
		}
	}
	return bytes
}
