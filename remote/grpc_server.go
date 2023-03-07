package remote

import (
	"context"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/protof"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcServer struct {
	protof.RpcServiceServer
	callbackFunc func(in *protof.RpcInfo) []byte
}

func (s *GrpcServer) CallbackResult(f func(in *protof.RpcInfo) []byte) {
	s.callbackFunc = f
}

func (s *GrpcServer) ListenAddr(addr string) {
	common.AssertNil(s.callbackFunc, "回调方法为空")
	// 监听连接
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("监听: %v", err)
	}
	v := grpc.NewServer()
	protof.RegisterRpcServiceServer(v, s)
	log.Println("服务端Rpc地址:" + addr)
	if err := v.Serve(lis); err != nil {
		log.Fatalf("监听失败: %v", err)
	}
}

// InvokeRemoteFunc 此处由Grpc客户端调用
func (s *GrpcServer) InvokeRemoteFunc(ctx context.Context, in *protof.RpcInfo) (*protof.RpcInfo, error) {
	common.AssertNil(s.callbackFunc, "没有注册回调方法")
	in.Body = s.callbackFunc(in)
	return in, nil
}

//func (s *GrpcServer) CountRoom() {
//	jsonDecoder := decoder.JsonDecoder{}
//	msg := message.JsonMessage{Merge: common.CmdKit.GetMerge(1, 2)}
//	var list []string
//	client := GrpcClient{}
//	ip := s.register.ListIp()
//	for _, info := range ip {
//		bytes := client.InvokeRemoteRpc(info.Ip+":"+fmt.Sprint(info.Port), msg.GetBytesResult())
//		result := jsonDecoder.DecoderBytes(bytes)
//		list = append(list, string(result.GetBody()))
//	}
//	log.Println(list)
//}
