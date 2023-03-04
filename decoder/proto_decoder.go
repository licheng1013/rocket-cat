package decoder

import (
	"github.com/licheng1013/io-game-go/message"
	"google.golang.org/protobuf/proto"
)

type ProtoDecoder struct {
}

func (p ProtoDecoder) EncodeBytes(result interface{}) []byte {
	bytes := message.MsgKit.StructToBytes(result.(proto.Message))
	return bytes
}

func (p ProtoDecoder) DecoderBytes(bytes []byte) message.Message {
	msg := message.ProtoMessage{}
	// 转换反序列话
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		panic(err)
	}
	return &msg
}
