package decoder

import (
	"core/message"
	"google.golang.org/protobuf/proto"
	"log"
)

type ProtoDecoder struct {
}

func (p ProtoDecoder) DecoderBytes(bytes []byte) (int64, interface{}) {
	msg := message.ProtoMessage{}
	// 转换反序列话
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		log.Println(err)
	}
	return msg.GetMerge(), msg.GetBody()
}
