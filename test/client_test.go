package main

import (
	"github.com/xtaci/kcp-go/v5"
	"io-game-go/message"
	"io-game-go/router"
	"testing"
)

func TestClient(t *testing.T) {
	kecClient, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defaultMessage := message.DefaultMessage{Body: "Hello", Merge: router.GetMerge(0, 1)}

	_, _ = kecClient.Write(message.GetObjectToBytes(defaultMessage))
	select {}
}
