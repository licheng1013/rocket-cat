package rpc

import "core/register"

// Rpc 远程调用接口,你可以随意实现自己的远程调用！
type Rpc interface {
	// Call 注册中心参数，路由，客户端的参数
	Call(requestInfo register.RequestInfo, merge int64, body interface{}) interface{}
}
