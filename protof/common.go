package protof

import (
	"google.golang.org/protobuf/proto"
	"log"
)

// RpcBodyMarshal 编码
func RpcBodyMarshal(v *RpcBody) []byte {
	body, err := proto.Marshal(v)
	if err != nil {
		log.Println("Proto编码错误:", err.Error())
	}
	return body
}

// RpcBodyUnmarshal 解码
func RpcBodyUnmarshal(body []byte, d *RpcBody) {
	err := proto.Unmarshal(body, d)
	if err != nil {
		log.Println("Proto解码错误:", err.Error())
	}
}

// RpcBodyBuild 构建一个
func RpcBodyBuild(body []byte) []byte {
	r := &RpcBody{Body: body}
	return RpcBodyMarshal(r)
}
