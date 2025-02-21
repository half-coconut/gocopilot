package book

// Queue 定义了一个简单的队列结构，先进先出
type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{
		items: make([]T, 0),
	}
}

// Make 添加一堆元素
func (q *Queue[T]) Make(items []T) {
	for _, v := range items {
		q.items = append(q.items, v)
	}
}

// Put 入队
func (q *Queue[T]) Put(item T) {
	q.items = append(q.items, item)
}

// Get 出队
func (q *Queue[T]) Get() T {
	if !q.IsEmpty() {
		first := q.items[0]
		q.items = q.items[1:]
		return first
	} else {
		panic("get from empty queue")
	}
}

// IsEmpty 判空
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}
