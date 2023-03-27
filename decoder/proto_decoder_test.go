package decoder

import (
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/protof"
	"testing"
)

type user2 struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func TestProtoDecoder(t *testing.T) {
	info := &protof.RpcInfo{Body: []byte("Hello")}
	// 增加简单方法
	decoder := ProtoDecoder{}
	data := decoder.Tool(1, 1, info)
	msg := decoder.DecoderBytes(data)
	info = &protof.RpcInfo{}
	_ = msg.Bind(info)
	t.Log(info)
	// 验证
	if string(info.Body) == "Hello" {
		t.Log("测试成功")
	} else {
		t.Log("测试失败")
	}

	// 测试
	message := &messages.ProtoMessage{}
	message.Headers = "HelloWorld"
	bytes := decoder.EncodeBytes(message)
	msg = decoder.DecoderBytes(bytes)
	t.Log(msg)
	if message.Headers == msg.GetHeaders() {
		t.Log("测试成功")
	} else {
		t.Log("测试失败")
	}
}
