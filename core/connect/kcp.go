// Package connect https://github.com/xtaci/kcp-go
package connect

import (
	"github.com/xtaci/kcp-go/v5"
	"log"
)

// ListenerKcp Kcp监听方法！
func ListenerKcp(addr string) {
	log.Println("服务器监听:" + addr)
	lis, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		log.Println("监听异常:", err)
	}
	for {
		conn, err := lis.AcceptKCP()
		if err != nil {
			log.Println("监听异常:", err)
		}
		go func(conn *kcp.UDPSession) {
			var buf = make([]byte, 4096)
			for {
				// 读取长度 n
				n, err := conn.Read(buf)
				if err != nil {
					log.Println(err)
					return
				}
				log.Println(string(buf[:n]))
				n, err = conn.Write(buf[:n])
				if err != nil {
					log.Println(err)
					return
				}
			}
		}(conn)
	}

}
