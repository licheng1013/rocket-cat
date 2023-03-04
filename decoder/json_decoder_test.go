package decoder

import (
	"github.com/licheng1013/io-game-go/message"
	"testing"
)

type user1 struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func TestJsonDecoder(t *testing.T) {
	u := user1{"小明", 12}
	// 问题
	jsonMessage := message.JsonMessage{Merge: 10, Code: -1, Message: "测试信息", Headers: "扩展参数", Heartbeat: true}
	jsonMessage.SetBody(message.MsgKit.StructToBytes(u))
	// 优化
	decoder := JsonDecoder{}
	msg := decoder.DecoderBytes(jsonMessage.GetBytesResult())
	t.Log(msg)
	var v user1
	err := jsonMessage.Bind(&v)
	if err != nil {
		panic(err)
	}
	t.Log(v)
}
