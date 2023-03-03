package common

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var listData = make(chan []byte)

// 定义一个打印函数，用于模拟一个耗时的任务
func printData() {
	//fmt.Println("ok") // 打印参数
	listData <- []byte("Hello")
	time.Sleep(time.Second)
}

func TestThread(t *testing.T) {

	go func() {
		for bytes := range listData {
			log.Println(string(bytes))
		}
	}()

	// 创建一个线程池，指定工作协程数为3，任务队列大小为10
	pool := NewPool(20, 30)
	// 启动线程池
	pool.Start()
	// 循环创建20个任务，并添加到线程池中，每个任务打印自己的编号和当前时间
	for i := 0; i < 20; i++ {
		err := pool.AddTask(printData, time.Second*3) // 向线程池中添加任务，指定超时时间为5秒
		if err != nil {                               // 如果添加失败，则打印错误信息并跳过该任务
			fmt.Println(err)
			continue
		}
	}
	// 停止线程池，并等待所有工作协程退出
	pool.Stop()
	close(listData)
}
