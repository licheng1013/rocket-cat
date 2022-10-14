package decoder

import (
	"io-game-go/message"
	"io-game-go/router"
)

var decoder Decoder = DefaultDecoder{}

// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte)
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) DecoderBytes(bytes []byte) {
	msg := message.GetBytesToObject(bytes)

	// TODO 这里先固定写死这个类型
	m := message.DefaultMessage{}

	router.GetObjectToToMap(msg, &m)
	router.ExecuteFunc(m.GetMerge(), m)
}

// GetDecoder 获取编码器
func GetDecoder() Decoder {
	return decoder
}

// SetDecoder 设置编码器
func SetDecoder(v Decoder) {
	decoder = v
}
