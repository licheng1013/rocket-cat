package room

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	// 测试匹配队列
	queue := NewMatchQueue(2, func(matches []IPlayer) {
		fmt.Println("匹配成功 -> ", matches)
	})
	queue.AddMatch(&DefaultPlayer{Uid: 1})
	queue.AddMatch(&DefaultPlayer{Uid: 2})
	queue.AddMatch(&DefaultPlayer{Uid: 3})
	queue.AddMatch(&DefaultPlayer{Uid: 4})
}
