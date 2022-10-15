package core

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
	"io"
	"io-game-go/decoder"
	"io-game-go/message"
	"io-game-go/router"
	"log"
	"net"
	"sync"
	"time"
)

var lock sync.Mutex
var v = 0
var t = time.Now().UnixMilli()

type Server struct {
	Port int
}

// SetDecoder 设置编码器默认: decoder.DefaultDecoder
func (g Server) SetDecoder(d decoder.Decoder) {
	decoder.SetDecoder(d)
}

func NewGameServer() *Server {
	g := &Server{}
	g.Port = 10000
	return g
}

func (g Server) Run() {
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	log.Println("监听端口: ", g.Port)

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
			go ConnTimer(conn)

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

				// 编码解码
				merge, body := decoder.GetDecoder().DecoderBytes(buffer[:n])
				// 处理对于函数
				result := router.ExecuteFunc(merge, body)
				if result != nil {
					// 分发消息
					var bytes []byte
					switch result.(type) {
					case []byte:
						bytes = result.([]byte)
						break
					case proto.Message:
						bytes = message.MarshalBytes(result.(proto.Message))
						break
					}
					_, err := conn.Write(bytes)
					if err != nil {
						log.Panicln(err)
					}
				}
				lock.Lock()
				v++
				lock.Unlock()

				startTime := time.Now().UnixMilli()
				if startTime-t > 1000 {
					log.Println("执行次数: ", v)
					v = 0
					t = startTime
				}

				//fmt.Println("receive from client:", buffer[:n])
			}
		}(conn)
	}
}
