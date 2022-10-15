package core

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"io-game-go/decoder"
	"io-game-go/router"
	"log"
	"net"
)

type GameServer struct {
	Port int
}

func NewGameServer() *GameServer {
	g := &GameServer{}
	g.Port = 10000
	return g
}

func (g GameServer) Run() {
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	log.Println("kcp listens on ", g.Port)

	addr := ":" + fmt.Sprint(g.Port)
	lis, err := kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		panic(err)
	}
	for {
		conn, e := lis.AcceptKCP()
		if e != nil {
			log.Panicln(e)
		}
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

				merge, body := decoder.GetDecoder().DecoderBytes(buffer[:n])
				result := router.ExecuteFunc(merge, body)
				if result != nil {
					// 分发消息
					switch result.(type) {
					case []byte:
						_, err := conn.Write(result.([]byte))
						if err != nil {
							log.Panicln(err)
						}
					}
				}
				//fmt.Println("receive from client:", buffer[:n])
			}
		}(conn)
	}
}
