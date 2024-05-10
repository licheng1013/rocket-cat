package common

import (
	"github.com/fatih/color"
	"log"
	"os"
)

// 颜色
var blueBg = color.New(color.FgBlue).SprintFunc()

// CatLog 一般日志
var CatLog *log.Logger

func init() {
	CatLog = log.New(os.Stdout, blueBg("[ROCKET CAT] "), log.LstdFlags)
}
