package common

import "github.com/panjf2000/ants/v2"

// Pool 定义一个线程池类型，包含任务队列和工作协程数
type Pool struct {
}

// NewPool 创建一个新的线程池
func NewPool() *Pool {
	// 当您调用此函数时，将预先存储池的全部容量
	pool := &Pool{}
	return pool
}

// AddTask 向线程池中添加一个新的任务，如果队列已满，则阻塞等待直到有空位或超时返回错误。
func (p *Pool) AddTask(task func()) {
	err := ants.Submit(task)
	if err != nil {
		RocketLog.Println("线程池错误 -> " + err.Error())
	}
}
