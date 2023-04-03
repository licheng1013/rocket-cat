package main

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/registers"
	"github.com/licheng1013/rocket-cat/router"
)

func main() {
	nacos := registers.DefaultNacos()
	// rpc
	service := core.DefaultService()
	service.SetRegister(nacos)
	service.Router().AddAction(1, 2, func(ctx *router.Context) {
		ctx.Data = []byte("1")
	})
	// 关机钩子
	service.AddClose(func() {
		common.Logger().Println("在关机中了")
	})
	service.Start()
}
