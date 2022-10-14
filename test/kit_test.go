package main

import (
	"io-game-go/message"
	"log"
	"testing"
)

func TestKit(t *testing.T) {
	defaultMessage := message.DefaultMessage{Body: "Hello", Merge: 2}
	bytes := message.GetObjectToBytes(defaultMessage)
	log.Println("转换为字节: ", bytes)
	object := message.GetBytesToObject(bytes)
	log.Println("转换为对象: ", object)
}
