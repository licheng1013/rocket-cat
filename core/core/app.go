package core

import (
	"core/common"
	"core/decoder"
	"core/pkc"
	"core/register"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// App 整个框架实例，请使用 NewGameServer 初始化
type App struct {
	// 监听端口
	port uint64
	// Kcp 连接监听 TODO 后续需要调整！
	listener *kcp.Listener
	// 启动前钩子
	beforeFunc func()
	// 开启收到消息日志
	EnableMessageLog bool
	// 解码器默认使用Json解码编码
	decoder decoder.Decoder
	// 关机钩子
	stopFunc func()
	// rpc请求
	rpc pkc.Rpc
	// 注册中心
	register register.Register
	// ip地址
	ip string
}

// SetDecoder 设置编码器
func (g *App) SetDecoder(d decoder.Decoder) {
	g.decoder = d
}

// NewGameServer 获取一个框架实例
func NewGameServer(register register.Register) *App {
	g := &App{}
	g.port = 10020
	g.beforeFunc = func() {}
	g.stopFunc = func() {}
	g.decoder = decoder.JsonDecoder{}
	g.rpc = pkc.DefaultRpc{}
	g.register = register
	return g
}

// SetProt 设置启动端口
func (g *App) SetProt(port uint64) {
	g.port = port
}

// Run 启动框架
func (g *App) Run() {
	g.beforeFunc()
	go func() {
		log.Println("监听端口: ", g.port)
		addr := ":" + fmt.Sprint(g.port)
		lis, err := kcp.ListenWithOptions(addr, nil, 10, 3)
		common.AssertErr(err)
		g.listener = lis
		for {
			conn, err := g.listener.AcceptKCP()
			common.AssertErr(err)
			go func(conn net.Conn) {
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
					merge, body := g.decoder.DecoderBytes(buffer[:n])
					if g.EnableMessageLog {
						log.Println("请求路由: ", merge, "请求数据: ", string(body.([]byte)))
					}

					rpcResult := &pkc.RpcResult{}
					// 处理对于函数 TODO 这里进行远程调用！
					err := g.rpc.Call(g.register.RequestUrl(), pkc.RequestInfo{Merage: merge, Body: body}, rpcResult)
					common.AssertErr(err)
					bytes := decoder.ParseResult(rpcResult.Result)
					if len(bytes) != 0 {
						_, err := conn.Write(bytes)
						common.AssertErr(err)
					}
				}
			}(conn)
		}
	}()

	log.Println("等待关闭...")
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("正在关机...")
	g.stopFunc()
	_ = g.listener.Close()
}

// SetStopFunc Stop 框架停止钩子！ 启动之前停止他
func (g *App) SetStopFunc(v func()) {
	g.stopFunc = v
}

// SetBeforeFunc 注册前置钩子，在框架启动的时候处理某些东西！
func (g *App) SetBeforeFunc(v func()) {
	g.beforeFunc = v
}

// SetIp 设置ip
func (g *App) SetIp(ip string) {
	g.ip = ip
}
