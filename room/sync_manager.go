package room

import (
	"time"
)

// FrameSyncManager 帧同步管理器
type FrameSyncManager struct {
	frameRate   float64       // 帧率
	delay       time.Duration // 延迟
	startTime   time.Time     // 开始时间
	currentTime time.Time     // 当前时间
}

// NewFrameSyncManager 创建新的帧同步管理器
func NewFrameSyncManager(frameRate float64, delay time.Duration) *FrameSyncManager {
	return &FrameSyncManager{
		frameRate: frameRate,
		delay:     delay,
	}
}

// Start 开始计时
func (f *FrameSyncManager) Start() {
	f.startTime = time.Now()
}

// WaitNextFrame 等待下一帧
func (f *FrameSyncManager) WaitNextFrame(callback func()) {
	f.currentTime = time.Now()
	// 帧率 = 1秒 / 帧率 * 当前帧
	duration := time.Duration(float64(time.Second)/f.frameRate) * time.Duration(f.GetCurrentFrame())
	expectedTime := f.startTime.Add(duration + f.delay)
	// 预期时间大于当前时间，就等待
	if expectedTime.After(f.currentTime) {
		// 等待预期时间减去当前时间则 = 休眠时间
		time.Sleep(expectedTime.Sub(f.currentTime))
	}
	callback()
}

// GetCurrentFrame 获取当前帧
func (f *FrameSyncManager) GetCurrentFrame() int {
	return int(f.currentTime.Sub(f.startTime).Seconds() * f.frameRate) // 当前时间减去开始时间 -> 秒，然后乘以帧率 = 当前帧
}
