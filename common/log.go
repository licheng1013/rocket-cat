package common

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger
var file *os.File
var lock sync.Mutex

// FileLogger 写入文件日志 -> 记录一些可能不重要的日志，例如客户端主动断开的错误。
func FileLogger() *log.Logger {
	lock.Lock()
	if file == nil {
		file, _ = os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	}
	if logger == nil {
		logger = log.New(file, "", log.LstdFlags+log.Lshortfile)
	}
	lock.Unlock()
	return logger
}
