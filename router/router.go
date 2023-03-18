package router

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"log"
	"time"
)

// Router 路由器功能
type Router interface {
	// AddAction 添加路由
	AddAction(cmd, subCmd int64, method func(ctx *Context))
	// ExecuteMethod 执行函数
	ExecuteMethod(msg *Context)
	// AddProxy 添加代理
	AddProxy(proxy Proxy)
}

// DefaultRouter 路由功能
type DefaultRouter struct {
	// 路由Id : 目标方法
	routerMap   map[int64]func(ctx *Context)
	middlewares Proxy
	DebugLog    bool
}

// AddAction 添加函数
func (r *DefaultRouter) AddAction(cmd, subCmd int64, method func(msg *Context)) {
	merge := common.CmdKit.GetMerge(cmd, subCmd)
	if r.routerMap == nil {
		r.routerMap = map[int64]func(msg *Context){}
	}
	if r.routerMap[merge] == nil {
		r.routerMap[merge] = method
		LogFunc(merge, method)
		return
	}
	panic(fmt.Sprintf("路由重复: %v-%v ", cmd, subCmd))
}

// InvokeFunc 代理函数执行
func (r *DefaultRouter) InvokeFunc(ctx *Context) {
	merge := ctx.Message.GetMerge()
	if r.routerMap[merge] == nil {
		log.Printf("路由: %v-%v 未注册", common.CmdKit.GetCmd(merge), common.CmdKit.GetSubCmd(merge))
		ctx.Message = nil
		return
	}
	if r.DebugLog {
		startTime := time.Now()
		r.routerMap[merge](ctx)
		entTime := time.Now()
		invokeTime := entTime.UnixMilli() - startTime.UnixMilli()
		LogFuncTime(merge, fmt.Sprint(invokeTime)+"ms")
	} else {
		r.routerMap[merge](ctx)
	}
}

func (r *DefaultRouter) SetProxy(proxy Proxy) {
	panic("被代理的类不允许设置代理对象!")
}

func (r *DefaultRouter) ExecuteMethod(msg *Context) {
	//var v Proxy
	//v = &ProxyFunc{r}
	//for i := range r.middlewares { //循环代理
	//	proxy := r.middlewares[i] //获取代理
	//	proxy.SetProxy(v)         //设置代理对象
	//	v = proxy                 //把当前对象设置下个代理
	//}
	if r.middlewares == nil {
		r.InvokeFunc(msg)
	} else {
		r.middlewares.InvokeFunc(msg)
	}
}

// AddProxy 添加代理器，最先添加的最最后执行。 异常中间件应该在最后添加机，用于捕获所有异常  ProxyFunc
func (r *DefaultRouter) AddProxy(proxy Proxy) {
	if r.middlewares == nil {
		r.middlewares = r
	}
	proxy.SetProxy(r.middlewares)
	r.middlewares = proxy
}
