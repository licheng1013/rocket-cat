package main

import (
	"core/message"
	"google.golang.org/protobuf/proto"
	"log"
	"testing"
)

func TestProtoDecoder(t *testing.T) {
	msg := message.ProtoMessage{Merge: 1}

	// 转换为proto
	marshal, err := proto.Marshal(&msg)
	if err != nil {
		panic(err)
	}
	log.Println(marshal)
	// 转换反序列话
	err = proto.Unmarshal(marshal, &msg)
	if err != nil {
		panic(err)
	}
	log.Println(msg.String())
}
