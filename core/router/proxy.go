package router

import (
	"github.com/io-game-go/message"
)

type Proxy interface {
	// InvokeFunc 调用函数
	InvokeFunc(msg message.Message) []byte
	// SetProxy 默认目标对象是A 你编写了B代理 -> 那么A将会传递下去,所以你在代理类需要调用代理对象的方法
	SetProxy(proxy Proxy)
}

// ProxyFunc 代理模型
type ProxyFunc struct {
	proxy Proxy
}

func (p *ProxyFunc) InvokeFunc(msg message.Message) []byte {
	return p.proxy.InvokeFunc(msg)
}

func (p *ProxyFunc) SetProxy(proxy Proxy) {
	p.proxy = proxy
}
