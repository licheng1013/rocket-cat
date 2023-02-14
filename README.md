# IoGameGo
## 介绍
- 2022/10/14
- 目前还是一个实验性项目

### 描述
- 一个go简单游戏服务器实现，目前是的。
- 打造一个简单的游戏服务器框架，易扩展，易使用。
- 网关与逻辑服通过Grpc进行调试
- 框架大部分功能都可以重写覆盖掉，自定义非常容易。

### 架构图
- 从架构图可以看出从 A获取B的地址是从通过注册中心**Nacos**获取的,当然你也可以自定义其他注册中心。
- ![struct.png](struct.png)

## 功能
### 传输结构
- [x] 支持Json
- [x] 支持Proto
- 当然你只要实现了此接口，就可以替换框架的实现了

```go
// Decoder 对数据的解码器
type Decoder interface {
	// DecoderBytes 收到客户端的数据
	DecoderBytes(bytes []byte) message.Message
	// EncodeBytes 封装编码
	EncodeBytes(result interface{}) []byte
}
```

### 连接相关
- [x] 支持Kcp
- [ ] 支持Tcp
- [x] 支持Websocket
- 如果你想自定义socket实现此接口即可。
- 上层并不关心内部如何处理的只需要把[]byte回写上去。

```go
type Socket interface {
    // ListenBack 监听连接收到的消息，回写到上层方法，当返回 byte 不为空时则写入到客户端
    ListenBack(func([]byte) []byte)
    ListenAddr(addr string)
}
```

### 负载均衡
- [x] Ncaos已提供


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