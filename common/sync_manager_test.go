package common

// 编写测试案例

import (
	"fmt"
	"testing"
	"time"
)

func TestSyncManager(t *testing.T) {
	frameRate := 60.0
	delay := time.Second / 60 // 以每秒60帧的速度运行，延迟60/1毫秒

	frameSyncManager := NewFrameSyncManager(frameRate, delay)
	frameSyncManager.Start()

	for i := 0; i < 60; i++ {
		frameSyncManager.WaitNextFrame(func() {
			fmt.Println("执行回调函数")
		})
		fmt.Printf("第 %d 帧在 %v\n", frameSyncManager.GetCurrentFrame(), time.Now())
	}
}

// 这个帧同步管理器的作用是确保程序以指定的帧率运行，以避免出现帧率不稳定的情况。它通过计算当前帧与下一帧之间的时间差来等待适当的时间，以确保下一帧在正确的时间到达。这对于需要精确控制帧率的应用程序非常有用，例如游戏或动画。
