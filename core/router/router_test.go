package router

import (
	"fmt"
	"github.com/io-game-go/message"
	"testing"
)

func TestRouter(t *testing.T) {
	router := DefaultRouter{}
	router.AddProxy(&B{})
	router.AddProxy(&C{})
	router.AddFunc(10, func(msg message.Message) []byte {
		fmt.Println("具体业务")
		return msg.GetBody()
	})
	_ = router.ExecuteFunc(&message.JsonMessage{Merge: 10})
}
