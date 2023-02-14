package pkc

import (
	"fmt"
	"github.com/io-game-go/message"
	"github.com/io-game-go/register"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Grpc struct {
	conn []*grpc.ClientConn
}

func (g *Grpc) Call(requestUrl registers.RequestInfo, info message.Message, rpcResult *RpcResult) error {
	if len(g.conn) <= 100 { //TODO 设置最大连接数！
		conn, err := grpc.Dial(fmt.Sprintf("%v:%v", requestUrl.Ip, requestUrl.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		g.conn = append(g.conn, conn)
	}

	return nil
}
