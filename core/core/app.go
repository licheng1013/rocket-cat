package core

import (
	"context"
	"core/decoder"
	"core/message"
	"core/router"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Port       uint64
	Listener   *kcp.Listener
	BeforeFunc func()
}

// SetDecoder 设置编码器默认: decoder.DefaultDecoder
func (g *App) SetDecoder(d decoder.Decoder) {
	decoder.SetDecoder(d)
}

func NewGameServer() *App {
	g := &App{}
	g.Port = 10000
	g.BeforeFunc = func() {}
	return g
}

func (g *App) Run() {
	g.BeforeFunc()
	go func() {
		log.Println("监听端口: ", g.Port)
		addr := ":" + fmt.Sprint(g.Port)
		lis, err := kcp.ListenWithOptions(addr, nil, 10, 3)
		AssertErr(err)
		g.Listener = lis
		for {
			conn, err := g.Listener.AcceptKCP()
			AssertErr(err)
			go func(conn net.Conn) {
				conn.RemoteAddr().String()
				var buffer = make([]byte, 1024, 1024)
				for {
					// 读取长度 n
					n, e := conn.Read(buffer)
					if e != nil {
						if e == io.EOF {
							break
						}
						fmt.Println(errorx.Wrap(e))
						break
					}
					// 编码解码
					merge, body := decoder.GetDecoder().DecoderBytes(buffer[:n])
					// 处理对于函数
					result := router.ExecuteFunc(merge, body)
					if result != nil {
						// 分发消息
						var bytes []byte
						switch result.(type) {
						case []byte:
							bytes = result.([]byte)
							break
						case proto.Message:
							bytes = message.MarshalBytes(result.(proto.Message))
							break
						}
						_, err := conn.Write(bytes)
						AssertErr(err)
					}
				}
			}(conn)
		}

	}()

}




// Stop 框架停止钩子！
func (g *App) Stop(v func()) {
	log.Println("等待关闭...")
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("正在关机...")
	// 创建一个5秒超时的context
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	v()
	_ = g.Listener.Close()
}

// SetBeforeFunc 注册前置钩子，在框架启动的时候处理某些东西！
func (g *App) SetBeforeFunc(v func()) {
	g.BeforeFunc = v
}
