package book

import "testing"

func TestStack(t *testing.T) {
	apple := []int{3, 2, 1, 6}
	s := NewStack[int]()
	s.Make(apple)
	println(s.Size())
	println(s.Pop())
	s.Push(5)
	println(s.Pop())
	println(s.Pop())
	println(s.Pop())
	println(s.Pop())
	println(s.Pop())
}

func TestQueue(t *testing.T) {
	apple := []int{3, 2, 1, 6}
	s := NewQueue[int]()
	s.Make(apple)
	println(s.Size())
	println(s.Get())
	s.Put(5)
	println(s.Get())
	println(s.Get())
	println(s.Get())
	println(s.Get())
	println(s.Get())
}
