package decoder

import "github.com/licheng1013/rocket-cat/messages"

// Decoder 对数据的解码器
type Decoder interface {
	//  收到客户端的数据
	Decoder(bytes []byte) messages.Message
	//  封装编码
	Encode(result any) []byte
	// Tool 工具方法
	Data(cmd, subCmd int64, body any) []byte
}
