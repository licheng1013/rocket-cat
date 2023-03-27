package router

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/messages"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	// 进行一亿次运算
	a := A{}
	c := ProxyFunc{proxy: &a}
	d := B{proxy: &c}
	e := ErrProxy{proxy: &d}
	for i := 0; i < 10; i++ {
		go func() {
			// 开始时间
			startTime := time.Now().UnixMilli()
			for k := 0; k < 10000000; k++ {
				e.InvokeFunc(&Context{Message: &messages.JsonMessage{}})
			}
			//结束时间
			endTime := time.Now().UnixMilli()
			fmt.Println("耗时:", endTime-startTime, "ms")
		}()
	}
	time.Sleep(10 * time.Second)
}

type A struct {
}

func (p *A) InvokeFunc(ctx *Context) {
	//FileLogger().Println("业务执行")
	//panic("系统错误") //测试异常
	//panic(NewServiceError(100, "业务异常"))
	ctx.Message.SetBody([]byte("ok"))
}
func (p *A) SetProxy(proxy Proxy) {

}

type B struct {
	proxy Proxy
}

func (p *B) InvokeFunc(ctx *Context) {
	//FileLogger().Println("B执行")
	p.proxy.InvokeFunc(ctx)
	//FileLogger().Println("B之后")
}
func (p *B) SetProxy(proxy Proxy) {
	p.proxy = proxy
}
