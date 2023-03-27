package decoder

import (
	"github.com/licheng1013/rocket-cat/messages"
	"testing"
)

func TestJsonDecoder(t *testing.T) {
	// 问题
	jsonMessage := messages.JsonMessage{Merge: 10, Code: -1, Message: "测试信息", Headers: "扩展参数", Heartbeat: true}
	// 优化
	decoder := JsonDecoder{}
	bytes := decoder.EncodeBytes(jsonMessage)
	msg := decoder.DecoderBytes(bytes)
	t.Log(msg)

	// 测试
	data := decoder.Tool(1, 1, 1)
	msg = decoder.DecoderBytes(data)
	t.Log(string(msg.GetBody()))
	if string(msg.GetBody()) == "1" {
		t.Log("测试成功")
	} else {
		t.Log("测试失败")
	}
}
