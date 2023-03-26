package common

import (
	"fmt"
	"testing"
	"time"
)

var count int64

// 定义一个打印函数，用于模拟一个耗时的任务
func printData() {
	//log.Println("ok") // 打印参数
	//listData <- []byte("Hello")
	//time.Sleep(time.Second)
	count++
}

func TestFunc(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100000000; i++ {
		printData()
	}
	// 结束时间
	end := time.Now()
	fmt.Println("耗时：", end.UnixMilli()-start.UnixMilli(), "ms", count)
}
