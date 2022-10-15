package decoder

import (
	"io-game-go/message"
	"io-game-go/router"
)

var decoder Decoder = DefaultDecoder{}

// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte) message.Message
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) DecoderBytes(bytes []byte) message.Message {
	msg := message.GetBytesToObject(bytes)

	// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
	m := message.DefaultMessage{}
	router.GetObjectToToMap(msg, &m)
	return m
}

// GetDecoder 获取编码器
func GetDecoder() Decoder {
	return decoder
}

// SetDecoder 设置编码器
func SetDecoder(v Decoder) {
	decoder = v
}
