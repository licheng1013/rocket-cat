package decoder

import (
	"core/message"
	"testing"
)

type user1 struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func TestJsonDecoder(t *testing.T) {
	u := user1{"小明", 12}
	// 问题
	jsonMessage := message.JsonMessage{}
	jsonMessage.SetBody(message.MsgKit.StructToBytes(u))
	// 优化
	decoder := JsonDecoder{}
	msg := decoder.DecoderBytes(jsonMessage.GetBytesResult())
	t.Log(msg)
	var v user1
	message.MsgKit.BytesToStruct(msg.GetBody(), &v)
	t.Log(v)
}
