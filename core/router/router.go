package router

import (
	"fmt"
	"log"
)

// Router 路由功能
type Router struct {
	// 路由Id : 目标方法
	RouterMap map[int64]func(msg interface{}) interface{}
}

// AddFunc 添加路由,如果需要返回给客户端数据，需要返回数组字节: []byte  默认的消息实现: message.DefaultMessage 是把里面的body返回到上传
func (r Router) AddFunc(merge int64, method func(msg interface{}) interface{}) {
	if r.RouterMap[merge] == nil {
		r.RouterMap[merge] = method
		return
	}
	panic(fmt.Sprintf("路由重复: %v-%v ", CmdKit.GetCmd(merge), CmdKit.GetSubCmd(merge)))
}

// ExecuteFunc 执行消息
func (r Router) ExecuteFunc(merge int64, v interface{}) interface{} {
	if r.RouterMap[merge] == nil {
		log.Printf("路由: %v-%v 未注册", CmdKit.GetCmd(merge), CmdKit.GetSubCmd(merge))
		return nil
	}
	return r.RouterMap[merge](v)
}
