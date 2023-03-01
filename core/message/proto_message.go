package message

import (
	"errors"
	"github.com/io-game-go/common"
	"github.com/io-game-go/protof"
	"google.golang.org/protobuf/proto"
)

type ProtoMessage struct {
	protof.Message
}

func (p *ProtoMessage) Bind(v interface{}) (err error) {
	if err = common.AssertPtr(v, "不是指针类型,无法绑定到结构体上"); err != nil {
		return
	}
	err = errors.New("不是 proto.Message 类型")
	switch v.(type) {
	case proto.Message:
		err = proto.Unmarshal(p.GetBody(), v.(proto.Message))
	}
	return
}

func (p *ProtoMessage) GetMerge() int64 {
	return p.Merge
}

func (p *ProtoMessage) GetBody() []byte {
	return p.Body
}

func (p *ProtoMessage) GetHeartbeat() bool {
	return p.Heartbeat
}

func (p *ProtoMessage) GetCode() int64 {
	return p.Code
}

func (p *ProtoMessage) GetMessage() string {
	return p.Message.Message
}

func (p *ProtoMessage) GetBytesResult() []byte {
	marshal, err := proto.Marshal(p)
	if err != nil {
		print(err)
	}
	return marshal
}

func (p *ProtoMessage) SetBody(bytes []byte) Message {
	p.Body = bytes
	return p
}

func (p *ProtoMessage) GetHeaders() string {
	return p.Headers
}
