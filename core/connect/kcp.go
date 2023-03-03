// Package connect https://github.com/xtaci/kcp-go
package connect

import (
	"github.com/io-game-go/common"
	"github.com/xtaci/kcp-go/v5"
)

type KcpSocket struct {
	MySocket
}

func (socket *KcpSocket) ListenBack(f func([]byte) []byte) {
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
		go func(conn *kcp.UDPSession) {

			socket.AsyncResult(func(bytes []byte) {
				_, err = conn.Write(bytes)
				if err != nil {
					// log.Println("写入错误:", err)
					common.FileLogger().Println("kcp写入错误: " + err.Error())
					_ = conn.Close()
				}
			})
			// 创建一个线程池，指定工作协程数为3，任务队列大小为10
			socket.pool = common.NewPool(20, 30)
			socket.pool.Start()

			var buf = make([]byte, 4096)
			for {
				// 读取长度 n
				n, err := conn.Read(buf)
				if err != nil {
					common.FileLogger().Println("kcp读取错误:", err.Error())
					break
				}
				socket.InvokeMethod(buf[:n])

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
		}(conn)
	}

}
