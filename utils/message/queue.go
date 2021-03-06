package message

import (
	"fmt"
	"sync"
)

// Queue 消息队列
type Queue struct {
	sync.Mutex             // 锁
	run        bool        // 是否已经运行
	current    chan []byte // 接收 byte数组的管道
}

// NewQueue 创建消息队列
func NewQueue(l int) *Queue {
	var q = &Queue{}
	q.run = false
	q.current = make(chan []byte, l)
	return q
}

// Push 推入数据
func (q *Queue) Push(b []byte) {
	q.current <- b
}

// Run 开始监听
func (q *Queue) Run() {
	q.Lock()
	if q.run {
		q.Unlock()
		return
	}
	q.run = true
	q.Unlock()
	go func() {
		for {
			q.Lock()
			data := <-q.current
			fmt.Println(string(data))
			q.Unlock()
		}
	}()
}
