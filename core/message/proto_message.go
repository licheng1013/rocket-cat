package message

import (
	"core/protof"
	"google.golang.org/protobuf/proto"
)

type ProtoMessage struct {
	protof.Message
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

func (p *ProtoMessage) SetBody(bytes []byte) {
	p.Body = bytes
}

func (p *ProtoMessage) GetHeaders() string {
	return p.Headers
}
