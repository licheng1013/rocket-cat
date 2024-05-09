package decoder

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/messages"
)

type JsonDecoder struct {
}

// EncodeBytes 编码为字节
func (d JsonDecoder) Encode(result interface{}) []byte {
	switch result.(type) {
	case []byte:
		return result.([]byte)
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		common.CatLog.Println("JsonDecoder -> 转换失败")
		return []byte{}
	}
	return bytes
}

// DecoderBytes 处理客户端返回的数据
func (d JsonDecoder) Decoder(bytes []byte) messages.Message {
	msg := messages.JsonMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		common.CatLog.Println("JsonDecoder -> 解析失败")
	}
	return &msg
}

// Tool 工具方法
func (d JsonDecoder) Data(cmd, subCmd int64, body interface{}) []byte {
	message := messages.JsonMessage{Merge: common.CmdKit.GetMerge(cmd, subCmd)}
	message.SetBody(body)
	return d.Encode(&message)
}

// JsonDecoderBytes 工具方法
func JsonDecoderBytes(bytes []byte) messages.Message {
	j := &JsonDecoder{}
	return j.Decoder(bytes)
}
