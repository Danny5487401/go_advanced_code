package _1_unblock_queue

type QueueInterface[T any] interface {
	Put(value T)
	Pop() (T, bool)
	Size() int
}

type Queue[T any] struct {
	elements []T
}

// Put 将数据放入队列尾部
func (q *Queue[T]) Put(value T) {
	q.elements = append(q.elements, value)
}

// Pop 从队列头部取出并从头部删除对应数据
func (q *Queue[T]) Pop() (T, bool) {
	var value T
	if len(q.elements) == 0 {
		return value, true
	}
	value = q.elements[0]
	q.elements = q.elements[1:]
	return value, len(q.elements) == 0
}

// Size 队列大小
func (q *Queue[T]) Size() int {
	return len(q.elements)
}
