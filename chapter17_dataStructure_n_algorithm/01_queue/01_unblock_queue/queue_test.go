package _1_unblock_queue

import (
	"fmt"
	"testing"
)

func TestIntDequeue(t *testing.T) {
	var qInterface QueueInterface[int]
	qInterface = new(Queue[int])
	qInterface.Put(1)
	qInterface.Put(2)
	qInterface.Put(3)
	fmt.Println(qInterface)
}

func TestStringDequeue(t *testing.T) {
	var qInterface QueueInterface[string]
	qInterface = new(Queue[string])
	qInterface.Put("a")
	qInterface.Put("b")
	qInterface.Put("c")
	fmt.Println(qInterface)
}
