package message

import (
	"gitee.com/licheng1013/go-util/common"
	"google.golang.org/protobuf/proto"
	"log"
)

func GetObjectToBytes(a any) []byte {
	return []byte(common.JsonUtil.JsonToStr(a))
}

func GetBytesToObject(bytes []byte) any {
	var v interface{}
	common.JsonUtil.JsonToMap(string(bytes), &v)
	return v
}

// UnmarshalBytes 字节转换为对象
func UnmarshalBytes(bytes []byte, info proto.Message) {
	err := proto.Unmarshal(bytes, info)
	if err != nil {
		log.Panicln(err)
	}
}

// UnmarshalInterface 字节转换为对象
func UnmarshalInterface(bytes interface{}, info proto.Message) {
	UnmarshalBytes(bytes.([]byte), info)
}

// MarshalBytes 转换为字节对象
func MarshalBytes(info proto.Message) []byte {
	marshal, err := proto.Marshal(info)
	if err != nil {
		log.Println(err)
	}
	return marshal
}
