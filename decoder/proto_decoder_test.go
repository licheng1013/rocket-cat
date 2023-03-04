package decoder

import (
	"github.com/io-game-go/message"
	"github.com/io-game-go/protof"
	"google.golang.org/protobuf/proto"
	"testing"
)

type user2 struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

func TestProtoDecoder(t *testing.T) {
	u := &protof.RpcInfo{Body: []byte("小明")}
	marshal, _ := proto.Marshal(u)
	// 问题
	protoMessage := message.ProtoMessage{}
	protoMessage.SetBody(marshal)
	// 优化
	decoder := ProtoDecoder{}
	msg := decoder.DecoderBytes(protoMessage.GetBytesResult())
	t.Log(msg)
	var v protof.RpcInfo
	err := protoMessage.Bind(&v)
	if err != nil {
		panic(err)
	}
	t.Log(string(v.Body))
}
