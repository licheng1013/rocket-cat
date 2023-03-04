package common

import (
	"errors"
	"fmt"
	"time"
)

// Pool 定义一个线程池类型，包含任务队列和工作协程数
type Pool struct {
	queue      chan func() // 任务队列，用于存放待执行的任务
	numWorkers int         // 工作协程数，用于控制并发数
}

// NewPool 创建一个新的线程池
func NewPool(numWorkers int, queueSize int) *Pool {
	return &Pool{
		queue:      make(chan func(), queueSize), // 初始化任务队列，指定队列大小
		numWorkers: numWorkers,                   // 设置工作协程数
	}
}

// Start 启动线程池，让工作协程开始从任务队列中取出并执行任务
func (p *Pool) Start() {
	for i := 0; i < p.numWorkers; i++ { // 循环创建工作协程，数量由numWorkers决定
		go p.worker() // 启动每个工作协程，传入编号i
	}
}

func (p *Pool) worker() {
	for task := range p.queue { // 循环从任务队列中取出任务，直到队列关闭或为空
		task() // 执行任务
	}
}

// AddTask 向线程池中添加一个新的任务，如果队列已满，则阻塞等待直到有空位或超时返回错误。
func (p *Pool) AddTask(task func(), timeout time.Duration) error {
	select {
	case p.queue <- task: // 尝试向队列中发送任务，如果成功则返回nil错误。
		return nil
	case <-time.After(timeout): // 如果超过指定的超时时间，则返回错误。
		return errors.New("timeout")
	}
}

// AddTaskNonBlocking 向线程池中添加一个新的任务，如果队列已满，则立即返回错误。
func (p *Pool) AddTaskNonBlocking(task func()) error {
	select {
	case p.queue <- task: // 尝试向队列中发送任务，如果成功则返回nil错误。
		return nil
	default: // 如果队列已满，则返回错误。
		return errors.New("queue full")
	}
}

// Stop 关闭线程池，停止接收新的任务，并等待所有工作协程完成当前的任务后退出。
func (p *Pool) Stop() {
	close(p.queue)         // 关闭任务队列，不再接收新的任务。
	for len(p.queue) > 0 { // 等待所有已经发送到队列中的任务被处理完毕。
		time.Sleep(time.Second)
	}
	fmt.Println("线程停止") // 打印线程池停止信息。
}
