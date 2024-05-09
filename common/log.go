package common

import (
	"github.com/fatih/color"
	"log"
	"os"
)

var loggerFile *log.Logger
var file *os.File

// FileLogger 写入文件日志 -> 记录一些可能不重要的日志，例如客户端主动断开的错误。
func FileLogger() *log.Logger {
	return loggerFile
}

var blueBg = color.New(color.FgBlue).SprintFunc()

// CatLog 一般日志
var CatLog *log.Logger

func init() {
	file, _ = os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	loggerFile = log.New(file, "", log.LstdFlags+log.Lshortfile)

	CatLog = log.New(os.Stdout, blueBg("[ROCKET CAT] "), log.LstdFlags)
}
