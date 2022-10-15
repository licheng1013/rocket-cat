package core

import (
	"log"
	"net"
	"time"
)

func ConnTimer(conn net.Conn) {
	for {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		_, err := conn.Write(make([]byte, 0))
		if err != nil {
			log.Println("关闭连接！")
			err := conn.Close()
			if err != nil {
				log.Panicln(err)
			}
		}
	}
}
