package common

import "sync"

type SafeList struct {
	mu    sync.Mutex
	items []any
}

// Add 添加元素
func (l *SafeList) Add(item any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, item)
}

// Len 获取长度
func (l *SafeList) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.items)
}

// Get 获取元素
func (l *SafeList) Get(index int) any {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.items[index]
}

// Remove 删除元素
func (l *SafeList) Remove(index int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items[:index], l.items[index+1:]...)
}

// GetList 获取列表
func (l *SafeList) GetList() []any {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.items
}
