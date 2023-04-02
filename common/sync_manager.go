package common

import (
	"time"
)

// FrameSyncManager 帧同步管理器
type FrameSyncManager struct {
	frameRate   float64       // 帧率
	frameDelay  time.Duration // 帧延迟
	delay       time.Duration // 延迟
	startTime   time.Time     // 开始时间
	currentTime time.Time     // 当前时间
}

// NewFrameSyncManager 创建新的帧同步管理器
func NewFrameSyncManager(frameRate float64, delay time.Duration) *FrameSyncManager {
	return &FrameSyncManager{
		frameRate:  frameRate,
		frameDelay: time.Duration(float64(time.Second) / frameRate), // 帧延迟
		delay:      delay,
	}
}

// Start 开始计时
func (f *FrameSyncManager) Start() {
	f.startTime = time.Now()
}

// WaitNextFrame 等待下一帧
//  1. 获取当前时间
//  2. 计算预期时间 = 开始时间 + 帧延迟 * 当前帧 + 延迟  -> 示例: 1秒 + 1/60秒 * 0 + 1/60秒 = 1.016秒
//  3. 如果预期时间大于当前时间，就等待; 示例: 1.016秒 > 1秒 就等待16毫秒
//  4. 执行回调函数
func (f *FrameSyncManager) WaitNextFrame(callback func()) {
	f.currentTime = time.Now()
	// 开始时间 + 帧延迟 * 当前帧 + 延迟 = 预期时间
	expectedTime := f.startTime.Add(time.Duration(float64(f.frameDelay)*float64(f.GetCurrentFrame())) + f.delay)
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
