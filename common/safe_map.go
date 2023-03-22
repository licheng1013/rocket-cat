package common

import (
	"sync"
)

// SafeMap 安全的map，支持并发读写
type SafeMap struct {
	sync.RWMutex
	m    map[int]interface{}
	keys []int
}

// NewSafeMap 创建一个安全的map
func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[int]interface{})}
}

// 设置值
func (sm *SafeMap) Set(key int, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := sm.m[key]; !ok {
		sm.keys = append(sm.keys, key)
	}
	sm.m[key] = value
}

// 获取值
func (sm *SafeMap) Get(key int) (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	value, ok := sm.m[key]
	return value, ok
}

// 获取所有键
func (sm *SafeMap) Keys() []int {
	sm.RLock()
	defer sm.RUnlock()
	return sm.keys
}

// 获取所有值
func (sm *SafeMap) Values() []interface{} {
	sm.RLock()
	defer sm.RUnlock()
	values := make([]interface{}, len(sm.keys))
	for i, key := range sm.keys {
		values[i] = sm.m[key]
	}
	return values
}
