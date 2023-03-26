// Package connect https://github.com/xtaci/kcp-go
package connect

import (
	"github.com/licheng1013/rocket-cat/common"
	"github.com/xtaci/kcp-go/v5"
)

type KcpSocket struct {
	MySocket
}

func (socket *KcpSocket) ListenBack(f func(uuid uint32, message []byte) []byte) {
	socket.proxyMethod = f
}

func (socket *KcpSocket) ListenAddr(addr string) {
	if socket.proxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	socket.listenerKcp(addr)
}

// listenerKcp Kcp监听方法！
func (socket *KcpSocket) listenerKcp(addr string) {
	//router.FileLogger().Println("服务器监听:" + addr)
	lis, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	for {
		conn, err := lis.AcceptKCP()
		if err != nil {
			common.FileLogger().Println("监听异常:", err.Error())
		}
		go socket.handleConn(conn)
	}

}

func (socket *KcpSocket) handleConn(conn *kcp.UDPSession) {
	uuid := socket.getNewChan()
	socket.AsyncResult(uuid, func(bytes []byte) {
		_, err := conn.Write(bytes)
		if socket.handleErr(err, uuid, "kcp写入错误: ") {
			return
		}
	})
	var buf = make([]byte, 4096)
	for {
		// 读取长度 n
		n, err := conn.Read(buf)
		if socket.handleErr(err, uuid, "kcp读取错误: ") {
			break
		}
		socket.InvokeMethod(uuid, buf[:n])
	}
}
