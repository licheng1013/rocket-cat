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
}
