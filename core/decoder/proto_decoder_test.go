package decoder

import (
	"github.com/io-game-go/message"
	"testing"
)

type user2 struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func TestProtoDecoder(t *testing.T) {
	u := user2{"小明", 12}
	// 问题
	protoMessage := message.ProtoMessage{}
	protoMessage.SetBody(message.MsgKit.StructToBytes(u))
	// 优化
	decoder := ProtoDecoder{}
	msg := decoder.DecoderBytes(protoMessage.GetBytesResult())
	t.Log(msg)
	var v user2
	message.MsgKit.BytesToStruct(msg.GetBody(), &v)
	t.Log(v)
}
