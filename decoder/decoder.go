package decoder

import (
	"io-game-go/message"
	"io-game-go/router"
)

var decoder Decoder = DefaultDecoder{}

// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte) (int64, interface{})
}

type DefaultDecoder struct {
}

// DecoderBytes 实现此方法可获取一些功能！
// 数据解析处理的核心方法，在这里你新建一个类去实现  Decoder 来完成自定义的的玩法
// 默认的编码器消息 message.DefaultMessage 需要进行json转换
func (d DefaultDecoder) DecoderBytes(bytes []byte) (int64, interface{}) {
	// 这里转换成了map
	msg := message.GetBytesToObject(bytes)
	// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
	m := message.DefaultMessage{}
	router.GetObjectToToMap(msg, &m)
	return m.GetMerge(), m.Body
}

// GetDecoder 获取编码器
func GetDecoder() Decoder {
	return decoder
}

// SetDecoder 设置编码器
func SetDecoder(v Decoder) {
	decoder = v
}
