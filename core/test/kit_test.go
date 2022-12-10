package main

import (
	"core/message"
	"core/router"
	"log"
	"testing"
)

func TestKit(t *testing.T) {
	defaultMessage := message.JsonMessage{Body: []byte("Hello"), Merge: 2}
	bytes := message.GetObjectToBytes(defaultMessage)
	log.Println("转换为字节: ", bytes)

	msg := message.GetBytesToObject(bytes)
	m := message.JsonMessage{}
	router.GetObjectToToMap(msg, &m)

	log.Println("转换为对象: ", m)
	log.Println("转换为对象: ", string(m.Body))
}
