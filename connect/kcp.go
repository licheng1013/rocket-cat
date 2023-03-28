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
	lis, err := kcp.ListenWithOptions(addr, nil, 0, 0)
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
	socketId := socket.getNewChan()
	socket.AsyncResult(socketId, func(bytes []byte) {
		//log.Println("写入数据->" + string(bytes))
		_, err := conn.Write(bytes)
		if socket.handleErr(err, socketId, "kcp写入错误: ") {
			return
		}
	})
	var buf = make([]byte, 4096)
	for {
		// 读取长度 n
		n, err := conn.Read(buf)
		if socket.handleErr(err, socketId, "kcp读取错误: ") {
			break
		}
		//log.Println("读取数据->" + string(buf[:n]))
		socket.InvokeMethod(socketId, buf[:n])
	}
}
