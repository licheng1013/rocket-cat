package router

import (
	"fmt"
	"github.com/licheng1013/rocket-cat/common"
	"log"
	"time"
)

type H map[string]any
type M map[string]bool

// Router 路由器功能
type Router interface {
	Action(cmd, subCmd int64, method func(ctx *Context))
	ExecuteMethod(msg *Context)
	AddProxy(proxy Proxy)
}

// DefaultRouter 路由功能
type DefaultRouter struct {
	routerMap  map[int64]func(ctx *Context) // 路由表
	middlewares Proxy 					  // 代理
	DebugLog   bool 						  // 是否开启日志
	SkipLogMap map[int64]bool 			  // 跳过日志
}

func (r *DefaultRouter) AddSkipLog(cmd, subCmd int64) {
	merge := common.CmdKit.GetMerge(cmd, subCmd)
	if r.SkipLogMap == nil {
		r.SkipLogMap = make(map[int64]bool)
	}
	r.SkipLogMap[merge] = true
}

func (r *DefaultRouter) Action(cmd, subCmd int64, method func(msg *Context)) {
	merge := common.CmdKit.GetMerge(cmd, subCmd)
	if r.routerMap == nil {
		r.routerMap = make(map[int64]func(msg *Context))
	}
	if r.routerMap[merge] != nil {
		panic(fmt.Sprintf("路由重复: %v-%v ", cmd, subCmd))
	}
	r.routerMap[merge] = method
	LogFunc(merge, method)
}

func (r *DefaultRouter) InvokeFunc(ctx *Context) {
	merge := ctx.Message.GetMerge()
	if r.routerMap[merge] == nil {
		log.Printf("路由: %v-%v 未注册", common.CmdKit.GetCmd(merge), common.CmdKit.GetSubCmd(merge))
		ctx.Message = nil
		return
	}
	if r.DebugLog && !r.SkipLogMap[merge] {
		startTime := time.Now()
		r.routerMap[merge](ctx)
		invokeTime := time.Since(startTime).Milliseconds()
		LogFuncTime(merge, fmt.Sprintf("%dms", invokeTime))
	} else {
		r.routerMap[merge](ctx)
	}
}

func (r *DefaultRouter) SetProxy(proxy Proxy) {
	panic("被代理的类不允许设置代理对象!")
}

func (r *DefaultRouter) ExecuteMethod(msg *Context) {
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
