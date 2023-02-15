package router

import "github.com/io-game-go/message"

type Proxy interface {
	InvokeFunc(msg message.Message) []byte
}

type ProxyFunc struct {
	proxy Proxy
}

func (p *ProxyFunc) InvokeFunc(msg message.Message) []byte {
	return p.proxy.InvokeFunc(msg)
}
