package decoder

import (
	"core/message"
	"core/router"
)

type JsonDecoder struct {
}

// DecoderBytes 实现此方法可获取一些功能！
// 数据解析处理的核心方法，在这里你新建一个类去实现  Decoder 来完成自定义的的玩法
// 默认的编码器消息 message.DefaultMessage 需要进行json转换
func (d JsonDecoder) DecoderBytes(bytes []byte) message.Message {
	// 这里转换成了map
	msg := message.GetBytesToObject(bytes)
	// TODO 这里是对数据处理实现部分，目前这个支持固定到字类
	m := message.JsonMessage{}
	router.GetObjectToToMap(msg, &m)
	return &m
}
