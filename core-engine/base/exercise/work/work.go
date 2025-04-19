package main

import (
	"sync"
	"time"
)

var result = make(chan []*Result, 20)

type Worker interface {
	Execute(result chan []*Result, s *subtask)
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	s := &subtask{
		began: time.Now(),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Execute(result, s)
			}
			p.wg.Done()
		}()
	}
	return &p
}

// Run 提交工作到工作池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有goroutine停止工作
func (p *Pool) Shutdown() {
	p.wg.Wait()
	close(p.work)
	close(result)
}
