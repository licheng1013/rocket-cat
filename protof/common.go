package protof

import (
	"github.com/licheng1013/rocket-cat/common"
	"google.golang.org/protobuf/proto"
)

// RpcBodyMarshal 编码
func RpcBodyMarshal(v *RpcInfo) []byte {
	body, err := proto.Marshal(v)
	if err != nil {
		common.RocketLog.Println("Proto编码错误:", err.Error())
	}
	return body
}

// RpcBodyUnmarshal 解码
func RpcBodyUnmarshal(body []byte, d *RpcInfo) {
	err := proto.Unmarshal(body, d)
	if err != nil {
		common.RocketLog.Println("Proto解码错误:", err.Error())
	}
}

// RpcBodyBuild 构建一个
func RpcBodyBuild(body []byte) *RpcInfo {
	return &RpcInfo{Body: body}
}
