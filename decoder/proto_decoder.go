package decoder

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/messages"
	"google.golang.org/protobuf/proto"
)

type ProtoDecoder struct {
}

func (p ProtoDecoder) EncodeBytes(result interface{}) []byte {
	switch result.(type) {
	case []byte:
		return result.([]byte)
	case proto.Message:
		bytes, err := proto.Marshal(result.(proto.Message))
		if err != nil {
			common.RocketLog.Println("ProtoDecoder -> 转换失败")
			break
		}
		return bytes
	}
	return []byte{}
}

// Tool 工具方法,用于简化编码,减少代码量
func (p ProtoDecoder) Tool(cmd, subCmd int64, body proto.Message) []byte {
	merge := common.CmdKit.GetMerge(cmd, subCmd)
	message := messages.ProtoMessage{}
	message.Merge = merge
	message.SetBody(body)
	return p.EncodeBytes(&message)
}

func (p ProtoDecoder) DecoderBytes(bytes []byte) messages.Message {
	msg := messages.ProtoMessage{}
	// 转换反序列话
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		common.RocketLog.Println("ProtoDecoder -> 解析失败")
	}
	return &msg
}

// ProtoDecoderBytes 工具方法
func ProtoDecoderBytes(bytes []byte) messages.Message {
	j := &ProtoDecoder{}
	return j.DecoderBytes(bytes)
}
