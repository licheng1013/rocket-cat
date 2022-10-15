package main

import (
	"io-game-go/message"
	"io-game-go/router"
	"log"
	"testing"
)

func TestKit(t *testing.T) {
	defaultMessage := message.DefaultMessage{Body: []byte("Hello"), Merge: 2}
	bytes := message.GetObjectToBytes(defaultMessage)
	log.Println("转换为字节: ", bytes)

	msg := message.GetBytesToObject(bytes)
	m := message.DefaultMessage{}
	router.GetObjectToToMap(msg, &m)

	log.Println("转换为对象: ", m)
	log.Println("转换为对象: ", string(m.Body))
}
