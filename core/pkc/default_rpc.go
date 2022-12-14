package pkc

import (
	"core/common"
	"core/message"
	"core/register"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type DefaultRpc struct {

}

func (DefaultRpc) Call(requestUrl register.RequestInfo, info message.Message, rpcResult *RpcResult) error {
	//log.Println("执行远程调用信息: ", requestUrl)
	cli, err := rpc.DialHTTP("tcp", requestUrl.Ip+":"+fmt.Sprint(requestUrl.Port))
	if err != nil {
		log.Println(err)
	}
	err = cli.Call("Result.Invok", info, &rpcResult)
	common.AssertErr(err)
	log.Println("远程结果:", string(rpcResult.Result))
	return nil
}

func ( DefaultRpc) RpcListen(ip string, port uint64)  {
	go func() {
		/*将服务对象进行注册*/
		err := rpc.Register(new(Result))
		if err != nil {
			err.Error()
		}
		rpc.HandleHTTP()
		/* 固定端口进行监听*/
		listen, err := net.Listen("tcp", ip+":"+fmt.Sprint(port))
		log.Println("Rpc监听地址: "+ ip +":"+fmt.Sprint(port))
		if err != nil {
			panic(err.Error())
		}
		_ = http.Serve(listen, nil)
	}()
}
