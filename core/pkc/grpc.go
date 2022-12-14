package pkc

import (
	"context"
	"core/message"
	"core/register"
	"flag"
	"fmt"
	"gitee.com/licheng1013/go-util/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

type Grpc struct {

}
func NewGreeterClient(cc grpc.ClientConnInterface) RpcHandle {
	return &GrpcResult{cc}
}

func (Grpc) Call(requestUrl register.RequestInfo, info message.Message, rpcResult *RpcResult) error {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v",requestUrl.Ip,requestUrl.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewGreeterClient(conn)

	// Contact the server and print out its response.
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = c.Invok(info,rpcResult)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println("请求结果:"+string(rpcResult.Result))
	return nil
}

// RpcListen 被调用逻辑处理
func ( Grpc) RpcListen(ip string, p uint64)  {
	port := flag.Int("port", int(p), "The server port")
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d",ip,*port))
		if err != nil {
			log.Fatalf("监听失败: %v", err)
		}
		s := grpc.NewServer()
		RegisterGreeterServer(s, GrpcResult{})
		log.Printf("监听端口 %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func RegisterGreeterServer(s grpc.ServiceRegistrar, srv RpcHandle) {
	s.RegisterService(&GreeterServicedesc, srv)
}

// GreeterServicedesc 是 Greeter 服务的 grpc.ServiceDesc。它仅适用于直接与 grpc.RegisterService 一起使用，不能进行自省或修改（即使是副本）
var GreeterServicedesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*RpcHandle)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Invok",
		},
	},
	Streams:  []grpc.StreamDesc{},
}

type GrpcResult struct {
	cc grpc.ClientConnInterface
}

func (g GrpcResult) Invok(rpcInfo message.Message, rpcResulet *RpcResult) error {
	protoMessage := message.ProtoMessage{}
	common.JsonUtil.MapTosStruct(rpcInfo,&protoMessage)

	//fmt.Println("被调用了!")
	err := g.cc.Invoke(context.Background(), "Invok",&protoMessage,&rpcResulet)
	if err != nil {
		log.Println(err)
	}
	return nil
}