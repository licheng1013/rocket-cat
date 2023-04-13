package router

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/messages"
	"testing"
)

func TestProxy(t *testing.T) {
	a := A{} // 最后执行的业务
	var d = B{ProxyFunc: ProxyFunc{Proxy: &ProxyFunc{Proxy: &a}}}
	var e = ErrProxy{ProxyFunc: ProxyFunc{Proxy: &d}}
	e.SetProxy(&d)
	e.InvokeFunc(&Context{Message: &messages.JsonMessage{}})
}

type A struct {
	ProxyFunc
}

func (p *A) InvokeFunc(ctx *Context) {
	fmt.Println("业务执行")
	ctx.Message.SetBody([]byte("ok"))
}

type B struct {
	ProxyFunc
}

func (p *B) InvokeFunc(ctx *Context) {
	fmt.Println("B执行")
	p.Proxy.InvokeFunc(ctx)
}
