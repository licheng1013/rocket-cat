package core

import (
	"core/common"
	"core/decoder"
	"core/message"
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
	// 客户端连接
	Conns []net.Conn
	// 请求消息
	Result message.Message
	// 建立连接时候触发的插件！
	Connecteds []Plugin
	// 处理业务逻辑之前的插件！
	ServiceBefores []Plugin
}

// AddConnectPlugin 添加插件,连接时插件,建立连接时会触发此插件！
func (g *App) AddConnectPlugin(pluginFunc Plugin) {
	g.Connecteds = append(g.Connecteds, pluginFunc)
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
func NewGameServer(register register.Register) *App {
	g := &App{}
	g.port = 10020
	g.beforeFunc = func() {}
	g.stopFunc = func() {}
	g.decoder = decoder.JsonDecoder{}
	g.rpc = &pkc.DefaultRpc{}
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
		log.Println(fmt.Sprintf("监听端口: %v:%v", g.ip, g.port))
		addr := ":" + fmt.Sprint(g.port)
		lis, err := kcp.ListenWithOptions(addr, nil, 10, 3)
		common.AssertErr(err)
		g.listener = lis
		// TODO 未来加入TCP,WEBSOCKET
		for { //监听链接！
			conn, err := g.listener.AcceptKCP()
			common.AssertErr(err)
			go listenerKcp(conn, g)
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

// Kcp监听方法！
func listenerKcp(conn net.Conn, g *App) {
	g.Conns = append(g.Conns, conn)
	for i := range g.Connecteds { //执行插件逻辑
		g.Connecteds[i].Invok(g)
	}
	var buffer = make([]byte, 1024, 1024)
	for {
		// 读取长度 n
		n, e := conn.Read(buffer)
		if e != nil {
			if e == io.EOF {
				continue
			}
			fmt.Println(errorx.Wrap(e))
		}
		result, err := g.handle(buffer[:n])
		if err != nil {
			log.Println(err)
			continue
		}
		if len(result.GetBytesResult()) == 0 {
			continue
		}
		//获取结果返回！
		_, err = conn.Write(result.GetBytesResult())
		common.AssertErr(err)
	}
}

func (g *App) handle(bytes []byte) (result message.Message, err error) {
	// 编码解码
	msg := g.decoder.DecoderBytes(bytes)
	if g.EnableMessageLog {
		log.Println("请求路由: ", msg.GetMerge(), "请求数据: ", string(msg.GetBody()))
	}
	g.Result = msg
	for i := range g.ServiceBefores { //执行调用远程时候执行的插件逻辑！
		g.ServiceBefores[i].Invok(g)
	}
	rpcResult := pkc.RpcResult{}
	// 处理对于函数 TODO 这里进行远程调用！
	err = g.rpc.Call(g.register.RequestUrl(), msg, &rpcResult)
	byteData := decoder.ParseResult(rpcResult.Result)
	msg.SetBody(byteData)
	return msg, err
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
