package decoder

import "github.com/licheng1013/io-game-go/messages"

// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte) messages.Message
	// EncodeBytes 封装编码
	EncodeBytes(result interface{}) []byte
}
