package message

import (
	"go-util/util"
)

func GetObjectToBytes(a any) []byte {
	return []byte(util.JsonToStr(a))
}

func GetBytesToObject(bytes []byte) any {
	v := GetMessage()
	util.JsonToObj(string(bytes), &v)
	return v
}
