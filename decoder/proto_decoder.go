package decoder

import (
	"github.com/licheng1013/io-game-go/messages"
	"google.golang.org/protobuf/proto"
)

type ProtoDecoder struct {
}

func (p ProtoDecoder) EncodeBytes(result interface{}) []byte {
	bytes := messages.MsgKit.StructToBytes(result.(proto.Message))
	return bytes
}

func (p ProtoDecoder) DecoderBytes(bytes []byte) messages.Message {
	msg := messages.ProtoMessage{}
	// 转换反序列话
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		panic(err)
	}
	return &msg
}

// ProtoDecoderBytes 工具方法
func ProtoDecoderBytes(bytes []byte) messages.Message {
	j := &ProtoDecoder{}
	return j.DecoderBytes(bytes)
}
