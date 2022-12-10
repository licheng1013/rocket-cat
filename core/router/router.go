package router

import (
	"encoding/json"
	"log"
)

var routerMap = make(map[int64]func(msg interface{}) interface{})

type CmdInfo struct {
	Cmd      int64
	SubCmd   int64
	CmdMerge int64
}

// GetMerge 获取路由
func GetMerge(cmd, subCmd int64) int64 {
	return (cmd << 16) + subCmd
}

// GetCmd 获取主命令
func GetCmd(merge int64) int64 {
	return merge >> 16
}

// GetSubCmd 获取子命令
func GetSubCmd(merge int64) int64 {
	return merge & 0xFFFF
}

// AddFunc 添加路由,如果需要返回给客户端数据，需要返回数组字节: []byte  默认的消息实现: message.DefaultMessage 是把里面的body返回到上传
func AddFunc(merge int64, method func(msg interface{}) interface{}) {
	if routerMap[merge] == nil {
		routerMap[merge] = method
	} else {
		log.Panicln("添加了重复的路由请检查: ", merge, " 主命令: ", GetCmd(merge), " 子命令: ", GetSubCmd(merge))
	}
}

// ExecuteFunc 执行消息
func ExecuteFunc(merge int64, v interface{}) interface{} {
	return routerMap[merge](v)
}

// GetObjectToToMap 转换map为结构体
func GetObjectToToMap(param interface{}, data interface{}) {
	str, _ := json.Marshal(param)
	_ = json.Unmarshal(str, data)
}
