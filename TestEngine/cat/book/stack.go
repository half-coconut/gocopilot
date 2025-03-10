package book

// Stake 后进先出
type Stack[T any] struct {
	items []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		items: make([]T, 0),
	}
}

// IsEmpty 判空
func (s *Stack[T]) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}

// Make 添加一堆元素
func (s *Stack[T]) Make(items []T) {
	for _, v := range items {
		s.items = append(s.items, v)
	}
}

// Push 添加一个元素
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop 弹出最后一个，并删除最后一个
func (s *Stack[T]) Pop() T {
	if !s.IsEmpty() {
		top := s.items[s.Size()-1]
		s.items = s.items[:s.Size()-1]
		return top
	} else {
		panic("pop from empty stack")
	}
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}
