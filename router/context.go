package router

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/remote"
	"log"
	"reflect"
	"runtime"
	"sync"
)

type Context struct {
	// 具体消息 -> 当这个消息为空时则不返回数据回去
	Message messages.Message
	// Rpc服务 -> 服务
	RpcServer remote.RpcServer
	// 具体消息 -> 此消息比 Message 更具有优先级返回
	Data []byte
	// 链接Id -> 连接建立时的唯一id
	SocketId uint32
}

var yellow = color.New(color.FgYellow).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()

var logger *log.Logger
var lock sync.Mutex

// FileLogger 写入文件日志 -> 记录一些可能不重要的日志，例如客户端主动断开的错误。
func FileLogger() *log.Logger {
	if logger == nil {
		lock.Lock()
		logger = log.Default()
		lock.Unlock()
	}
	return logger
}

type routerInfo struct {
	merge string
	name  string
}

var infoMap = make(map[int64]*routerInfo, 0)

func LogFuncTime(merge int64, time string) {
	info := infoMap[merge]
	if info == nil {
		return
	}
	LogPrint("MERGE: %s | %s --> %s.\n", yellow(info.merge), time, blue(info.name))
}

func LogFunc(merge int64, f func(ctx *Context)) {
	info := infoMap[merge]
	if info == nil {
		cmd := common.CmdKit.GetCmd(merge)
		subCmd := common.CmdKit.GetSubCmd(merge)
		mergeInfo := fmt.Sprint(cmd) + "-" + fmt.Sprint(subCmd)
		info = &routerInfo{
			merge: mergeInfo,
			name:  runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
		}
		infoMap[merge] = info
	}
	LogPrint("MERGE: %s --> %s.\n", yellow(info.merge), blue(info.name))
}

func LogPrint(format string, values ...any) {
	FileLogger().Printf("[ROCKET CAT] "+format, values...)
}
