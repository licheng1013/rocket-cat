package router

import (
	"github.com/io-game-go/message"
	"log"
	"testing"
)

func TestProxy(t *testing.T) {
	// a - >
	a := A{}
	b := ProxyFunc{proxy: &a}
	c := ProxyFunc{proxy: &b}
	d := B{proxy: &c}
	d.InvokeFunc(&message.JsonMessage{})
}

type A struct {
}

func (p *A) InvokeFunc(msg message.Message) []byte {
	log.Println("业务执行")
	return make([]byte, 0)
}
func (p *A) SetProxy(proxy Proxy) {
}

type B struct {
	proxy Proxy
}

func (p *B) InvokeFunc(msg message.Message) []byte {
	log.Println("B执行")
	invokeFunc := p.proxy.InvokeFunc(msg)
	log.Println("B之后")
	return invokeFunc
}
func (p *B) SetProxy(proxy Proxy) {
	p.proxy = proxy
}

type C struct {
	proxy Proxy
}

func (p *C) InvokeFunc(msg message.Message) []byte {
	log.Println("C执行")
	invokeFunc := p.proxy.InvokeFunc(msg)
	log.Println("C之后")
	return invokeFunc
}
func (p *C) SetProxy(proxy Proxy) {
	p.proxy = proxy
}
