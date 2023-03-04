// Package connect https://github.com/xtaci/kcp-go
package connect

import (
	"github.com/licheng1013/io-game-go/common"
	"github.com/xtaci/kcp-go/v5"
)

type KcpSocket struct {
	MySocket
}

func (socket *KcpSocket) ListenBack(f func(uuid uint32, message []byte) []byte) {
	socket.proxyMethod = f
}

func (socket *KcpSocket) ListenAddr(addr string) {
	socket.init()
	if socket.proxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	socket.listenerKcp(addr)
}

// listenerKcp Kcp监听方法！
func (socket *KcpSocket) listenerKcp(addr string) {
	//log.Println("服务器监听:" + addr)
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

	uuid := common.UuidKit.UUID()
	messageChannel := make(chan []byte)
	socket.UuidOnCoon.Store(uuid, messageChannel)
	go func() {
		for data := range messageChannel {
			socket.queue <- data
		}
	}()

	socket.AsyncResult(func(bytes []byte) {
		_, err := conn.Write(bytes)
		if err != nil {
			// log.Println("写入错误:", err)
			common.FileLogger().Println("kcp写入错误: " + err.Error())
			_ = conn.Close()
			socket.close(uuid)
		}
	})

	var buf = make([]byte, 4096)
	for {
		// 读取长度 n
		n, err := conn.Read(buf)
		if err != nil {
			common.FileLogger().Println("kcp读取错误:", err.Error())
			socket.close(uuid)
			break
		}
		socket.InvokeMethod(uuid, buf[:n])

		// log.Printf("收到消息: %s", buf[:n])
		//bytes := socket.proxyMethod(buf[:n])
		//if len(bytes) == 0 {
		//	continue
		//}
		//n, err = conn.Write(bytes)
		//if err != nil {
		//	common.FileLogger().Println("kcp写入错误:", err.Error())
		//	break
		//}
	}
}
