package queue

/*
队列
	先进先出
实现
	用slice 展现,线程安全
主要方法：
	New()
	Enqueue()
	Dequeue()
	Front()
	IsEmpty()
	Size()
*/
import (
	"sync"

	"github.com/cheekybits/genny/generic"
)

// 队列中的内容
type Item generic.Type //使用泛型

type ItemQueue struct {
	items []Item
	lock  sync.RWMutex
}

// 生成队列
func (s *ItemQueue) New() *ItemQueue {
	s.items = []Item{}
	return s
}

//入队
func (s *ItemQueue) Enqueue(t Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, t)
}

//出队  从头部
func (s *ItemQueue) Dequeue() *Item {
	s.lock.Lock()
	defer s.lock.Unlock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	return &item
}

// 获取第一个元素，未清除
func (s *ItemQueue) Front() *Item {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return &item
}

// 判断是否清除
func (s *ItemQueue) IsEmpty() bool {
	return len(s.items) == 0
}

// 返回元素个数
func (s *ItemQueue) Size() int {
	return len(s.items)
}

/*
命令:
	genny -in queue.go -out queue-int.go gen "Item=int"
*/
