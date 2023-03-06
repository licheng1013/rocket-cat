package protof

import (
	"log"
	"testing"
)

func TestBody(t *testing.T) {
	v := &RpcBody{Body: []byte("Hi"), SocketId: 222}
	marshal := RpcBodyMarshal(v)
	d := &RpcBody{}
	RpcBodyUnmarshal(marshal, d)
	log.Println(d.SocketId)
	log.Println(string(d.Body))
}
