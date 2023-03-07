# RocketCat

![Alt](https://repobeats.axiom.co/api/embed/6e9456520132509e9335fb6ee214abacae172845.svg "Repobeats analytics image")

## 介绍

<p align="center">
<img align="center" width="150" src="images/cat-6047457_640.png">
</p>

- 2022/10/14
- 目前还是一个实验性项目
- 还需要很多要完成的东西
- **注意:此文档某些代码的Api已经变化了**

### 起步

- 安装后复制下列代码即可运行.
- go get github.com/licheng1013/rocket-cat

```go
package main

import (
	"github.com/gorilla/websocket"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/connect"
	"github.com/licheng1013/rocket-cat/core"
	"github.com/licheng1013/rocket-cat/decoder"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/router"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	socket := &connect.WebSocket{}
	channel := make(chan int)
	gateway := core.NewGateway()
	gateway.SetDecoder(decoder.JsonDecoder{})
	gateway.Router().AddAction(common.CmdKit.GetMerge(1, 1), func(ctx *router.Context) {
		socket.SendMessage(ctx.Message.SetBody([]byte("Hi")).GetBytesResult())
		ctx.Message.SetBody([]byte("Hi Ok 2"))
	})
	go gateway.Start(connect.Addr, socket)
	time.Sleep(time.Second / 2) //等待完全启动
	go WsTest(channel)
	select {
	case ok := <-channel:
		log.Println(ok)
	}
}

func WsTest(v chan int) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: connect.Addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	var count int64
	go func() {
		defer close(done)
		for {
			_, m, err := c.ReadMessage()
			jsonDecoder := decoder.JsonDecoder{}
			dto := jsonDecoder.DecoderBytes(m)
			if err != nil {
				log.Println("读取消息错误:", err)
				return
			}
			log.Println("收到消息:", string(dto.GetBody()))
			count++
			if count >= 2 {
				v <- 0
			}
		}
	}()
	for {
		jsonMessage := messages.JsonMessage{Body: []byte("HelloWorld")}
		jsonMessage.Merge = common.CmdKit.GetMerge(1, 1)
		err := c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
		if err != nil {
			log.Println("写:", err)
			return
		}
	}
}

```

### 描述

- 一个go简单游戏服务器实现，目前是的。
- 打造一个简单的游戏服务器框架，易扩展，易使用。
- 网关与逻辑服通过Grpc进行调用
- 框架大部分功能都可以重写覆盖掉，自定义非常容易。

### 架构图

- 从架构图可以看出从A获取B的地址是从通过注册中心**Nacos**获取的,当然你也可以自定义其他注册中心。
- 注:单机模式-不需要nacos即可使用
- ![struct.png](struct.png)

## 功能

- [x] 广播支持websocket.
- [x] 自定义底层连接.
- [x] 支持中间件.
- [x] 单机模式.
- [x] 自定义远程通信.
- [x] 网关支持kcp,websocket,tcp.
- [x] 负载均衡-由注册中心提供.
- [x] 逻辑服互调功能.
- [x] 传输协议json,proto


### 路由代理

- 在调用目标方法之前或之后处理一些自定义逻辑,需要实现 Proxy 接口
- 符合标准的实现 SetProxy 传入了目标代理对象 InvokeFunc 是执行目标方法。

```go
// ProxyFunc 代理模型
type ProxyFunc struct {
proxy Proxy
}

func (p *ProxyFunc) InvokeFunc(ctx Context) []byte {
return p.proxy.InvokeFunc(ctx)
}

func (p *ProxyFunc) SetProxy(proxy Proxy) {
p.proxy = proxy
}
```
