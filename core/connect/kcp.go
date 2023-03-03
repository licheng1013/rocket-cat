// Package connect https://github.com/xtaci/kcp-go
package connect

import (
	"github.com/io-game-go/common"
	"github.com/xtaci/kcp-go/v5"
)

type KcpSocket struct {
	MySocket
}

func (k *KcpSocket) ListenBack(f func([]byte) []byte) {
	k.proxyMethod = f
}

func (k *KcpSocket) ListenAddr(addr string) {
	if k.proxyMethod == nil {
		panic("未注册回调函数: ListenBack")
	}
	k.listenerKcp(addr)
}

// listenerKcp Kcp监听方法！
func (k *KcpSocket) listenerKcp(addr string) {
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
			var buf = make([]byte, 4096)
			for {
				// 读取长度 n
				n, err := conn.Read(buf)
				if err != nil {
					common.FileLogger().Println("kcp读取错误:", err.Error())
					break
				}
				// log.Printf("收到消息: %s", buf[:n])
				bytes := k.proxyMethod(buf[:n])
				if len(bytes) == 0 {
					continue
				}
				n, err = conn.Write(bytes)
				if err != nil {
					common.FileLogger().Println("kcp写入错误:", err.Error())
					break
				}
			}
		}(conn)
	}

}
