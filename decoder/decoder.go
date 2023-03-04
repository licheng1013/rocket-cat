package decoder

import (
	"github.com/io-game-go/message"
)

// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte) message.Message
	// EncodeBytes 封装编码
	EncodeBytes(result interface{}) []byte
}
