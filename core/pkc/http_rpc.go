package pkc

import (
	"core/register"
	"fmt"
	"log"
	"net/rpc"
)

type HttpRpc struct {
}

func (HttpRpc) Call(requestUrl register.RequestInfo, info RequestInfo, rpcResult *RpcResult) error {
	log.Println("执行远程调用信息: ", requestUrl)
	cli, err := rpc.DialHTTP("tcp", requestUrl.Ip+":"+fmt.Sprint(requestUrl.Port))
	if err != nil {
		panic(err)
	}
	err = cli.Call("Result.Invok", info, &rpcResult)
	if err != nil {
		panic(err)
	}
	fmt.Println("远程结果:", string(rpcResult.Result))
	return nil
}

type RequestInfo struct {
	Merage int64
	Body   interface{}
}

type RpcResult struct {
	Result []byte
}
