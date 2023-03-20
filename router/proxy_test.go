package router

import (
	"github.com/licheng1013/rocket-cat/messages"
	"testing"
)

func TestProxy(t *testing.T) {
	// a - >
	a := A{}
	c := ProxyFunc{proxy: &a}
	d := B{proxy: &c}
	d.InvokeFunc(&Context{Message: &messages.JsonMessage{}})
}

type A struct {
}

func (p *A) InvokeFunc(ctx *Context) {
	FileLogger().Println("业务执行")
	ctx.Message.SetBody([]byte("ok"))
}
func (p *A) SetProxy(proxy Proxy) {

}

type B struct {
	proxy Proxy
}

func (p *B) InvokeFunc(ctx *Context) {
	FileLogger().Println("B执行")
	p.proxy.InvokeFunc(ctx)
	FileLogger().Println("B之后")
}
func (p *B) SetProxy(proxy Proxy) {
	p.proxy = proxy
}
