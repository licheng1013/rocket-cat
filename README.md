# IoGameGo
## 介绍
- 2022/10/14
- 目前还是一个实验性项目

## 描述
- 一个go简单游戏服务器实现，目前是的。
- 打造一个简单的游戏服务器框架，易扩展，易使用。
- 网关与逻辑服通过Grpc进行调试

## 架构图
- ![struct.png](struct.png)

## 功能
### 传输结构
- [x] 支持Json
- [x] 支持Proto

### 连接协议
- [x] 支持Kcp
- [ ] 支持Tcp
- [ ] 支持Websocket

### 负载均衡
- [x] Ncaos已提供

## 注册中心  
### Nacos
- 单机启动： startup.cmd -m standalone
- [Nacos](https://nacos.io/zh-cn/docs/v2/quickstart/quick-start.html)

## 工具类
- [Lancet](https://github.com/duke-git/lancet/blob/main/README_zh-CN.md)

## 如何调试
- 运行 nacos 单机
- 运行 gateway 网关
- 运行 service 服务
- 运行 example 客户端

## 计划
- [x] 心跳功能预览完成
- [] 数据互相通信

## 示例
### 单机示例
- 2023/2/14
- 完成初步的消息处理解析
- gateway_test.go 此文件展示了一个基本demo

```go
package core

import (
	"core/connect"
	"core/decoder"
	"core/message"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	gateway := NewGateway()
	gateway.SetDecoder(decoder.JsonDecoder{})

	gateway.Router().AddFunc(10, func(msg message.Message) []byte {
		log.Println(string(msg.GetBody()))
		return msg.GetBody()
	})

	gateway.Start(connect.Addr, &connect.WebSocket{})
}

func TestWsClient(t *testing.T) {
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
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			jsonMessage := message.JsonMessage{Body: []byte(t.String())}
			jsonMessage.Merge = 10
			err := c.WriteMessage(websocket.TextMessage, jsonMessage.GetBytesResult())
			if err != nil {
				log.Println("写:", err)
				return
			}
		case <-interrupt:
			log.Println("中断")
			// 通过发送关闭消息干净地关闭连接，然后等待（超时）服务器关闭连接。
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("写关闭:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
```