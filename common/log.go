package common

import (
	"github.com/fatih/color"
	"log"
	"os"
	"sync"
)

var loggerFile *log.Logger
var file *os.File
var lock sync.Mutex

// FileLogger 写入文件日志 -> 记录一些可能不重要的日志，例如客户端主动断开的错误。
func FileLogger() *log.Logger {
	lock.Lock()
	if file == nil {
		file, _ = os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	}
	if loggerFile == nil {
		loggerFile = log.New(file, "", log.LstdFlags+log.Lshortfile)
	}
	lock.Unlock()
	return loggerFile
}

var blueBg = color.New(color.FgBlue).SprintFunc()

// RocketLog 一般日志
var RocketLog *log.Logger

func init() {
	RocketLog = log.New(os.Stdout, blueBg("[ROCKET CAT] "), log.LstdFlags)
}
