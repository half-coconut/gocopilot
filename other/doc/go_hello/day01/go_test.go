package day01

import (
	"fmt"
	"testing"
	"time"
)

// Slice
func TestSlice(t *testing.T) {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := slice[2:5]
	t.Log(s1)
	s2 := s1[2:6:7]
	t.Log(s2)

	s2 = append(s2, 100)
	s2 = append(s2, 200)

	s1[2] = 20
	t.Log(s1)
	t.Log(s2)
	t.Log(slice)
}

func TestSliceV2(t *testing.T) {
	s := []int{1, 1, 1}
	f(s)
	t.Log(s)
}

func f(s []int) {
	for _, i := range s {
		i++
	}
	for i := range s {
		s[i] += 1 // 这里已经改变了切片共享的底层数据结构，因此输入结果改变了
	}
}

// Channel

func TestChannel(t *testing.T) {
	ch := make(chan int)
	GoroutineA(ch)
	GoroutineB(ch)
	ch <- 3
	time.Sleep(time.Second)
	// 死锁
}

func GoroutineA(a <-chan int) {
	val := <-a
	fmt.Println(val)
	return
}
func GoroutineB(b <-chan int) {
	val := <-b
	fmt.Println(val)
	return
}
