package router

import (
	"encoding/json"
	"log"
)

var routerMap = make(map[int]func(msg interface{}) interface{})

type CmdInfo struct {
	Cmd      int
	SubCmd   int
	CmdMerge int
}

// GetMerge 获取路由
func GetMerge(cmd, subCmd int) int {
	return (cmd << 16) + subCmd
}

// GetCmd 获取主命令
func GetCmd(merge int) int {
	return merge >> 16
}

// GetSubCmd 获取子命令
func GetSubCmd(merge int) int {
	return merge & 0xFFFF
}

// AddFunc 添加路由,如果需要返回给客户端数据，需要返回数组字节: []byte
func AddFunc(merge int, method func(msg interface{}) interface{}) {
	if routerMap[merge] == nil {
		routerMap[merge] = method
	} else {
		log.Panicln("添加了重复的路由请检查: ", merge, " 主命令: ", GetCmd(merge), " 子命令: ", GetSubCmd(merge))
	}
}

// ExecuteFunc 执行消息
func ExecuteFunc(merge int, v interface{}) interface{} {
	return routerMap[merge](v)
}

// GetObjectToToMap 转换map为结构体
func GetObjectToToMap(param interface{}, data interface{}) {
	str, _ := json.Marshal(param)
	_ = json.Unmarshal(str, data)
}
