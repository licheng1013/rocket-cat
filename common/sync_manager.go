package common

import (
	"time"
)

// 帧同步管理器
type FrameSyncManager struct {
	frameRate   float64       // 帧率
	frameDelay  time.Duration // 帧延迟
	delay       time.Duration // 延迟
	startTime   time.Time     // 开始时间
	currentTime time.Time     // 当前时间
}

// 创建新的帧同步管理器
func NewFrameSyncManager(frameRate float64, delay time.Duration) *FrameSyncManager {
	return &FrameSyncManager{
		frameRate:  frameRate,
		frameDelay: time.Duration(float64(time.Second) / frameRate),
		delay:      delay,
	}
}

// 开始计时
func (f *FrameSyncManager) Start() {
	f.startTime = time.Now()
}

// 等待下一帧
func (f *FrameSyncManager) WaitNextFrame(callback func()) {
	f.currentTime = time.Now()
	expectedTime := f.startTime.Add(time.Duration(float64(f.frameDelay)*float64(f.GetCurrentFrame())) + f.delay)
	if expectedTime.After(f.currentTime) {
		time.Sleep(expectedTime.Sub(f.currentTime))
	}
	callback()
}

// 获取当前帧
func (f *FrameSyncManager) GetCurrentFrame() int {
	return int(f.currentTime.Sub(f.startTime).Seconds() * f.frameRate)
}
