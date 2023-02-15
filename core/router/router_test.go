package router

import (
	"fmt"
	"github.com/io-game-go/message"
	"testing"
)

func TestRouter(t *testing.T) {
	router := DefaultRouter{}
	router.AddProxy(&B{})
	router.AddFunc(10, func(ctx Context) []byte {
		fmt.Println("具体业务")
		return ctx.Message.GetBody()
	})
	_ = router.ExecuteFunc(Context{Message: &message.JsonMessage{Merge: 10}})
}
