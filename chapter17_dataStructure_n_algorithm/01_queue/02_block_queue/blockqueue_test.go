package _2_block_queue

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewBlockingQueue(t *testing.T) {
	a := NewBlockingQueue[int](context.Background(), 3)
	a.SetRejectHandler(func(ctx context.Context) bool {
		fmt.Println("reject")
		return false
	})
	for i := 0; i < 4; i++ {
		go a.Push(i)
	}
	time.Sleep(time.Second * 5)

	t.Log(a.Poll())
	t.Log(a.TryPoll(time.Second))

}
