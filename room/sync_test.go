package room

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
