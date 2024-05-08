package router

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/messages"
	"testing"
)

func TestRouter(t *testing.T) {
	msg := "HelloWorld"
	router := DefaultRouter{}
	router.DebugLog = true
	router.AddProxy(&B{})
	router.AddSkipLog(1, 2)
	merge := common.CmdKit.GetMerge(1, 2)
	router.Action(1, 2, func(ctx *Context) {
		fmt.Println("具体业务")
		fmt.Println("收到消息:" + string(ctx.Message.GetBody()))
		ctx.Message = nil
	})
	c := &Context{Message: &messages.JsonMessage{Merge: merge, Body: []byte(msg)}}
	router.ExecuteMethod(c)
	fmt.Println(c.Message)
}
