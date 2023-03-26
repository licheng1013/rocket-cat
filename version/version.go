package version

import (
	"github.com/fatih/color"
	"log"
	"os"
	"runtime"
)

const Version = "0.0.23"

func StartLogo() {
	// 获取go版本号
	goVersion := runtime.Version()
	// 设置颜色
	var blueBg = color.New(color.FgBlue).SprintFunc()
	var green = color.New(color.FgHiGreen).SprintFunc()
	// 打印logo
	log.New(os.Stderr, blueBg("[ROCKET CAT] "), log.LstdFlags).Println(
		green("\n" +
			"      /\\_/\\" +
			"\n     / o o \\" +
			"\n    =(   W  )=" +
			"\n     )     (" +
			"\n    (__\\_/__) Version -> " + Version + " Go Version -> " + goVersion))
}
