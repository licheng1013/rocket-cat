package decoder

import (
	"core/message"
	"google.golang.org/protobuf/proto"
	"log"
)

type ProtoDecoder struct {
}

func NewProtoDecoder() *ProtoDecoder {
	return &ProtoDecoder{}
}

func (p ProtoDecoder) DecoderBytes(bytes []byte) message.Message {
	msg := message.ProtoMessage{}
	// 转换反序列话
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		log.Println(err)
	}
	return &msg
}
