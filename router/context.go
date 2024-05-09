package router

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/licheng1013/rocket-cat/common"
	"github.com/licheng1013/rocket-cat/messages"
	"github.com/licheng1013/rocket-cat/remote"
	"log"
	"os"
	"reflect"
	"runtime"
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
	// 网关服的RpcIp,单机模式此处为空。将会传递到逻辑服当中
	RpcIp string
}

var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgHiGreen).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var blueBg = color.New(color.FgBlue).SprintFunc()

var logger *log.Logger

var greenText = green("MERGE")

// FileLogger 写入文件日志 -> 记录一些可能不重要的日志，例如客户端主动断开的错误。
func FileLogger() *log.Logger {
	return logger
}

func init() {
	logger = log.New(os.Stdout, blueBg("[ROCKET CAT] "), log.LstdFlags)
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
	LogPrint(greenText+": %s | %s \n", yellow(info.merge), time)
}

func LogFunc(merge int64, f func(ctx *Context)) {
	info := infoMap[merge]
	if info == nil {
		cmd := common.CmdKit.GetCmd(merge)
		subCmd := common.CmdKit.GetSubCmd(merge)
		mergeInfo := fmt.Sprint(cmd) + "-" + fmt.Sprint(subCmd) + " -> " + fmt.Sprint(merge)
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		info = &routerInfo{
			merge: mergeInfo,
			name:  name,
		}
		infoMap[merge] = info
	}
	LogPrint(greenText+": %s --> %s\n", yellow(info.merge), blue(info.name))
}

func LogPrint(format string, values ...any) {
	FileLogger().Printf(format, values...)
}
