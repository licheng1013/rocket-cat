# Rocket

- 不稳定.

## Start

- go get github.com/licheng1013/rocket-cat

```go
package main

import (
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/router"
)

func main() {
	// 构建一个默认服务
	gateway := core.DefaultGateway()
	// 添加一个路由
	gateway.Router().AddAction(1, 1, func(ctx *router.Context) {
		// 设置返回数据
		ctx.Message.SetBody(map[string]interface{}{"name": "RocketCat"})
		ctx.Message.SetMessage("这是消息")
	})
	// 绑定路由
	gateway.Start(":10100")
}

```
