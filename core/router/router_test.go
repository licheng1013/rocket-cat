package router

import (
	"fmt"
	"github.com/io-game-go/message"
	"testing"
)

func RouterTest(t *testing.T) {
	router := DefaultRouter{}
	router.AddFunc(10, func(msg message.Message) []byte {
		fmt.Println("具体业务")
		return msg.GetBody()
	})
}
