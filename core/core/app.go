package core

import (
	"github.com/io-game-go/decoder"
	"github.com/io-game-go/pkc"
	"github.com/io-game-go/register"
	"github.com/xtaci/kcp-go/v5"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	register registers.Register
	// ip地址
	ip string
	// 客户端连接
	Conns []net.Conn
	// 开启心跳功能
	EnableHearbeat bool
	// 心跳超时Map管理
	TimeOutMap sync.Map
}

// SetDecoder 设置编码器
func (g *App) SetDecoder(d decoder.Decoder) {
	g.decoder = d
}

// SetRpc 设置Rpc调用处理
func (g *App) SetRpc(p pkc.Rpc) {
	g.rpc = p
}

// NewGameServer 获取一个框架实例
func NewGameServer(register registers.Register) *App {
	g := &App{}
	g.port = 10020
	g.beforeFunc = func() {}
	g.stopFunc = func() {}
	g.decoder = decoder.JsonDecoder{}
	g.rpc = &pkc.DefaultRpc{}
	g.register = register
	g.TimeOutMap = sync.Map{}
	return g
}

// SetProt 设置启动端口
func (g *App) SetProt(port uint64) {
	g.port = port
}

// Run 启动框架
func (g *App) Run() {
	// 心跳检查器
	go func() {
		if !g.EnableHearbeat {
			return
		}
		log.Println("启动心跳功能...")
		for true {
			time.Sleep(3 * time.Second)
			g.TimeOutMap.Range(func(key, value any) bool {
				// 最大超时3秒
				if time.Now().UnixMilli()-value.(int64) > 3000 {
					log.Println("超时:", key.(uint32))
					for i := range g.Conns {
						session := g.Conns[i].(*kcp.UDPSession)
						if session.GetConv() == key.(uint32) {
							err := session.Close()
							if err != nil {
								log.Println("关闭连接:", key)
								g.TimeOutMap.Delete(key)
							}
							g.Conns = append(g.Conns[:i], g.Conns[i+1:]...) //删除这个元素
							break
						}
					}
				}
				return true
			})
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
