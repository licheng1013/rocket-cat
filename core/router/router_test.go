package router

import (
	"fmt"
	"github.com/io-game-go/message"
	"testing"
)

func TestRouter(t *testing.T) {
	msg := "HelloWorld"
	router := DefaultRouter{}
	router.AddProxy(&B{})
	merge := CmdKit.GetMerge(1, 2)
	router.AddFunc(merge, func(ctx Context) []byte {
		fmt.Println("具体业务")
		fmt.Println("收到消息:" + string(ctx.Message.GetBody()))
		return ctx.Message.GetBody()
	})
	_ = router.ExecuteMethod(Context{Message: &message.JsonMessage{Merge: merge, Body: []byte(msg)}})
}
