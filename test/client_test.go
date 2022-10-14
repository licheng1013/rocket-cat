package main

import (
	"github.com/xtaci/kcp-go/v5"
	"io-game-go/message"
	"io-game-go/router"
	"testing"
)

func TestClient1(t *testing.T) {
	kecClient, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.DefaultMessage{Body: "Hello", Merge: router.GetMerge(0, 1)}

	for i := 0; i < 100; i++ {
		go func() {
			for true {
				_, _ = kecClient.Write(message.GetObjectToBytes(defaultMessage))
			}
		}()
	}
	select {}
}

func TestClient2(t *testing.T) {
	kecClient, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.DefaultMessage{Body: "Hello", Merge: router.GetMerge(0, 1)}

	for i := 0; i < 100; i++ {
		go func() {
			for true {
				_, _ = kecClient.Write(message.GetObjectToBytes(defaultMessage))
			}
		}()
	}
	select {}
}
