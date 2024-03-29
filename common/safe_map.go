package common

import (
	"sync"
)

// SafeMap 安全的map，支持并发读写,加有序的key
type SafeMap struct {
	sync.RWMutex
	m    map[int64]interface{}
	keys []int64
}

// NewSafeMap 创建一个安全的map
func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[int64]interface{})}
}

// Set 设置值
func (sm *SafeMap) Set(key int64, value interface{}) {
	sm.Lock()
	defer sm.Unlock()
	if _, ok := sm.m[key]; !ok {
		sm.keys = append(sm.keys, key)
	}
	sm.m[key] = value
}

// Get 获取值
func (sm *SafeMap) Get(key int64) (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	value, ok := sm.m[key]
	return value, ok
}

// Keys 获取所有键
func (sm *SafeMap) Keys() []int64 {
	sm.RLock()
	defer sm.RUnlock()
	return sm.keys
}

// Values 获取所有值
func (sm *SafeMap) Values() []interface{} {
	sm.RLock()
	defer sm.RUnlock()
	values := make([]interface{}, len(sm.keys))
	for i, key := range sm.keys {
		values[i] = sm.m[key]
	}
	return values
}

// Delete 删除值
func (sm *SafeMap) Delete(key int64) {
	sm.Lock()
	defer sm.Unlock()
	delete(sm.m, key)
	for i, k := range sm.keys {
		if k == key {
			sm.keys = append(sm.keys[:i], sm.keys[i+1:]...)
			break
		}
	}
}