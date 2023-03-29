package common

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	// 测试匹配队列
	queue := NewMatchQueue(2, func(matches []int64) {
		fmt.Println("匹配成功 -> ", matches)
	})
	queue.AddMatch(1)
	queue.AddMatch(1)
	queue.AddMatch(5)
	queue.AddMatch(4)
}
