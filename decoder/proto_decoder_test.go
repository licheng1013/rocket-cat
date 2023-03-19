package decoder

import (
	"github.com/licheng1013/rocket-cat/messages"
	"testing"
)

type user2 struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func TestProtoDecoder(t *testing.T) {
	// 问题
	protoMessage := messages.ProtoMessage{}
	protoMessage.Code = -1
	protoMessage.Body = []byte("测试消息")
	protoMessage.Merge = 10
	protoMessage.Heartbeat = true
	protoMessage.Headers = "扩展参数"
	// 优化
	decoder := ProtoDecoder{}
	bytes := decoder.EncodeBytes(&protoMessage)
	msg := decoder.DecoderBytes(bytes)
	t.Log(msg)
}
