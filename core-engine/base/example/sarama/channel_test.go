package sarama

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
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
			// 写入
			ch <- i
		}
		close(ch)
	}()

	for val := range ch {
		// 读取
		t.Log("发送成功", val)
	}
	t.Log("channel 被关了")
}

func TestChannelBlock(t *testing.T) {
	ch := make(chan int)
	val := <-ch
	t.Log(val)
}

func TestGorutineCh(t *testing.T) {
	ch := make(chan int, 100000)
	// 这1个 就泄露了
	go func() {
		for i := 0; i < 100000; i++ {
			// 写入
			ch <- i
		}
		abc := new(BigObj)
		t.Log(abc)
		// 永久阻塞在这里，ch 占据的内存，永远不会被回收
		ch <- 1
	}()
	// 这里后面没有人往 ch 里面读数据
	t.Log(fmt.Sprintf("当前 goroutine 数量: %d\n", runtime.NumGoroutine()))
}

type BigObj struct {
}

func TestSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	go func() {
		time.Sleep(time.Millisecond * 100)
		ch1 <- 123
	}()
	go func() {
		time.Sleep(time.Millisecond * 100)
		ch2 <- 456
	}()
	select {
	case val := <-ch1:
		t.Log("ch1", val)
		val = <-ch2
		t.Log("ch2", val)
	case val := <-ch2:
		t.Log("ch2", val)
		val = <-ch1
		t.Log("ch1", val)
	}
	// 多个分支同时满足要求，随机执行一个
	// 如果所有case 都阻塞了，就执行 default
}
