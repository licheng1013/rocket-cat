package rpc

import (
	"core/register"
	"log"
)

type HttpRpc struct {
}

func (HttpRpc) Call(requestUrl register.RequestInfo, merge int64, body interface{}) interface{} {
	log.Println("执行远程调用!")
	return nil
}
