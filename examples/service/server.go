package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"time"
)

var sum int32

//func myFunc(i interface{}) {
//	n := i.(int32)
//	atomic.AddInt32(&sum, n)
//	fmt.Printf("run with %d\n", n)
//}

func demoFunc() {
	//time.Sleep(10 * time.Millisecond)
	//fmt.Println("Hello World!")
	sum++
}

func main() {
	defer ants.Release()
	runTimes := 10000000
	// 开始时间
	start := time.Now()
	for i := 0; i < runTimes; i++ {
		_ = ants.Submit(demoFunc)
	}
	//fmt.Printf("运行goroutines: %d\n", ants.Running())
	fmt.Printf("完成所有任务 -> %v\n", sum)
	// 结束时间并打印
	end := time.Now()
	fmt.Println("耗时：", end.Sub(start), "ms")
}
