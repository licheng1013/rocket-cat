package main

import (
	"core/pkc"
	"core/protof"
	"core/register"
	"fmt"
	"time"
)

func main() {
	ip := "192.168.101.10"
	prot := 8999

	// 获取服务单的消息
	unix := time.Now().UnixMilli()
	fmt.Println(unix)
	var count int64
	p := &pkc.Grpc{}
	result := pkc.RpcResult{}

	for true {
		count++
		err := p.Call(register.RequestInfo{Ip: ip, Port: uint64(prot)}, &protof.ProtoMessage{}, &result)
		if err != nil {
			panic(err)
		}
		//log.Println(string(result.Result))
		newUnix := time.Now().UnixMilli()
		if newUnix-unix > 1000 {
			fmt.Println("1秒请求数:", fmt.Sprint(count))
			unix = newUnix
			count = 0
		}
	}
	select {

	}
}
