package router

import (
	"core/message"
	"fmt"
	"log"
)

// Router 路由器功能
type Router interface {
	// AddFunc 添加路由
	AddFunc(merge int64, method func(msg message.Message) []byte)
	// InvokeFunc 执行目标方法之前
	InvokeFunc(msg message.Message) []byte
}

// DefaultRouter 路由功能
type DefaultRouter struct {
	// 路由Id : 目标方法
	routerMap map[int64]func(msg message.Message) []byte
}

// AddFunc 添加函数
func (r *DefaultRouter) AddFunc(merge int64, method func(msg message.Message) []byte) {
	if r.routerMap == nil {
		r.routerMap = map[int64]func(msg message.Message) []byte{}
	}
	if r.routerMap[merge] == nil {
		r.routerMap[merge] = method
		return
	}
	panic(fmt.Sprintf("路由重复: %v-%v ", CmdKit.GetCmd(merge), CmdKit.GetSubCmd(merge)))
}

// InvokeFunc 执行函数
func (r *DefaultRouter) InvokeFunc(msg message.Message) []byte {
	merge := msg.GetMerge()
	if r.routerMap[merge] == nil {
		log.Printf("路由: %v-%v 未注册", CmdKit.GetCmd(merge), CmdKit.GetSubCmd(merge))
		return nil
	}
	return r.routerMap[merge](msg)
}
