package router

import (
	"fmt"
	"log"
)

// Router 路由器功能
type Router interface {
	// AddFunc 添加路由
	AddFunc(merge int64, method func(ctx Context) []byte)
	// ExecuteFunc 执行函数
	ExecuteFunc(msg Context) []byte
}

// DefaultRouter 路由功能
type DefaultRouter struct {
	// 路由Id : 目标方法
	routerMap   map[int64]func(ctx Context) []byte
	middlewares []Proxy
}

// AddFunc 添加函数
func (r *DefaultRouter) AddFunc(merge int64, method func(msg Context) []byte) {
	if r.routerMap == nil {
		r.routerMap = map[int64]func(msg Context) []byte{}
	}
	if r.routerMap[merge] == nil {
		r.routerMap[merge] = method
		return
	}
	panic(fmt.Sprintf("路由重复: %v-%v ", CmdKit.GetCmd(merge), CmdKit.GetSubCmd(merge)))
}

// InvokeFunc 代理函数执行
func (r *DefaultRouter) InvokeFunc(ctx Context) []byte {
	merge := ctx.Message.GetMerge()
	if r.routerMap[merge] == nil {
		log.Printf("路由: %v-%v 未注册", CmdKit.GetCmd(merge), CmdKit.GetSubCmd(merge))
		return nil
	}
	return r.routerMap[merge](ctx)
}

func (r *DefaultRouter) SetProxy(proxy Proxy) {
	panic("被代理的类不允许设置代理对象!")
}

func (r *DefaultRouter) ExecuteFunc(msg Context) []byte {
	var v Proxy
	v = &ProxyFunc{r}
	for i := range r.middlewares {
		proxy := r.middlewares[i]
		proxy.SetProxy(v)
		v = proxy
	}
	return v.InvokeFunc(msg)
}

// AddProxy 添加代理器，最先添加的最最后执行。 异常中间件应该在最后天机，用于捕获所有异常  ProxyFunc
func (r *DefaultRouter) AddProxy(proxy Proxy) {
	r.middlewares = append(r.middlewares, proxy)
}
