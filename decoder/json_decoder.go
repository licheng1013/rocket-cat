package decoder

import (
	"encoding/json"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/messages"
)

type JsonDecoder struct {
}

// EncodeBytes 编码为字节
func (d JsonDecoder) EncodeBytes(result interface{}) []byte {
	switch result.(type) {
	case []byte:
		return result.([]byte)
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		common.Logger().Println("JsonDecoder -> 转换失败")
		return []byte{}
	}
	return bytes
}

// DecoderBytes 处理客户端返回的数据
func (d JsonDecoder) DecoderBytes(bytes []byte) messages.Message {
	msg := messages.JsonMessage{}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		common.Logger().Println("JsonDecoder -> 解析失败")
	}
	return &msg
}

// JsonDecoderBytes 工具方法
func JsonDecoderBytes(bytes []byte) messages.Message {
	j := &JsonDecoder{}
	return j.DecoderBytes(bytes)
}
