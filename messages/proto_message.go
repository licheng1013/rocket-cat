package messages

import (
	"errors"
	"github.com/licheng1013/io-game-go/common"
	"github.com/licheng1013/io-game-go/protof"
	"google.golang.org/protobuf/proto"
	"log"
)

type ProtoMessage struct {
	protof.Message
}

func (p *ProtoMessage) SetMerge(merge int64) {
	p.Merge = merge
}

func (p *ProtoMessage) SetCode(code int64) {
	p.Code = code
}

func (p *ProtoMessage) SetMessage(message string) {
	p.Message.Message = message
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

func (p *ProtoMessage) SetBody(data interface{}) Message {
	switch data.(type) {
	case []byte:
		p.Body = data.([]byte)
		break
	case proto.Message:
		marshal, err := proto.Marshal(data.(proto.Message))
		if err != nil {
			log.Println("Proto转换器错误,具体错误: ", err.Error())
		}
		p.Body = marshal
		break
	}
	return p
}

func (p *ProtoMessage) GetHeaders() string {
	return p.Headers
}
