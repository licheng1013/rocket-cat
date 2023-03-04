package message

import (
	"encoding/json"
)

var MsgKit = kit{}

type kit struct {
}

// StructToBytes 结构体转换为Json
func (k kit) StructToBytes(a any) (bytes []byte) {
	var err error
	if bytes, err = json.Marshal(a); err != nil {
		panic(err)
	}
	return
}

// BytesToStruct 字节转换为结构体
func (k kit) BytesToStruct(bytes []byte, v interface{}) (err error) {
	err = json.Unmarshal(bytes, &v)
	return
}
