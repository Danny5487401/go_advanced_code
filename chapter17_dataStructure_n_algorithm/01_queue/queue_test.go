package queue

import "testing"

var s ItemQueue

func initQueue() *ItemQueue {
	if s.items == nil {
		s = ItemQueue{}
		s.New()
	}
	return &s
}
func TestEnqueue(t *testing.T) {
	s := initQueue()
	s.Enqueue(1)
	s.Enqueue(2)
	s.Enqueue(3)
	if size := s.Size(); size != 3 {
		t.Errorf("wrong count, expected 3 and got %d", size)
	}
}

func TestDequeue(t *testing.T) {
	s.Enqueue(1)
	s.Enqueue(1)
	s.Enqueue(1)
	s.Dequeue()
	if size := len(s.items); size != 2 {
		t.Errorf("wrong count, expected 2 and got %d", size)
	}

	s.Dequeue()
	s.Dequeue()
	if size := len(s.items); size != 0 {
		t.Errorf("wrong count, expected 0 and got %d", size)
	}

	if !s.IsEmpty() {
		t.Errorf("IsEmpty should return true")
	}
}
