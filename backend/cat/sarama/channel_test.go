package sarama

import (
	"sync"
	"testing"
)

// channel
func TestChannel(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 123
	val := <-ch
	println(val)
}

func TestChannelClose(t *testing.T) {
	// closed 之后还可以读
	ch := make(chan int, 2)
	ch <- 123
	ch <- 234
	val, ok := <-ch
	t.Log(val, ok)
	close(ch)
	//ch <- 124
	val, ok = <-ch
	t.Log(val, ok)
	val, ok = <-ch
	t.Log(val, ok)
}

type MyStruct struct {
	ch chan struct{}
	// 一个实例只会执行一次
	closeOnce sync.Once
}

func (m *MyStruct) Close() {
	// 确保整个代码只会执行一次
	m.closeOnce.Do(func() {
		close(m.ch)
	})
}

func TestLoopChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	for val := range ch {
		t.Log("发送成功", val)
	}
	t.Log("channel 被关了")
}
