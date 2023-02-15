package router

import (
	"github.com/io-game-go/message"
	"log"
	"testing"
)

func TestProxy(t *testing.T) {
	// a - >
	a := A{}
	c := ProxyFunc{proxy: &a}
	d := B{proxy: &c}
	d.InvokeFunc(Context{Message: &message.JsonMessage{}})
}

type A struct {
}

func (p *A) InvokeFunc(ctx Context) []byte {
	log.Println("业务执行")
	return make([]byte, 0)
}
func (p *A) SetProxy(proxy Proxy) {

}

type B struct {
	proxy Proxy
}

func (p *B) InvokeFunc(ctx Context) []byte {
	log.Println("B执行")
	invokeFunc := p.proxy.InvokeFunc(ctx)
	log.Println("B之后")
	return invokeFunc
}
func (p *B) SetProxy(proxy Proxy) {
	p.proxy = proxy
}
