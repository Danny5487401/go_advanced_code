package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

/*
需求：
	计算大量整数和的程序
*/
type Task struct {
	index int   //组
	nums  []int // 需要计算的数
	sum   int   //总和
	wg    *sync.WaitGroup
}

func (t *Task) Do() {
	for _, num := range t.nums {
		t.sum += num
	}

	t.wg.Done()
}

// 需要执行的逻辑
func taskFunc(data interface{}) {
	task := data.(*Task)
	task.Do()
	fmt.Printf("task:%d sum:%d\n", task.index, task.sum)
}

const (
	DataSize    = 10000
	DataPerTask = 100
)

func main() {
	// ants.NewPoolWithFunc(cap, func(interface{}))这种方式创建的池子对象需要指定池函数，并且使用p.Invoke(arg)调用池函数。arg就是传给池函数func(interface{})的参数
	// 1.生成池对象
	p, _ := ants.NewPoolWithFunc(10, taskFunc)
	defer p.Release()

	// 模拟数据，做数据切分，生成任务
	// 随机生成 10000 个整数，将这些整数分为 100 份，每份 100 个
	nums := make([]int, DataSize, DataSize)
	rand.Seed(time.Now().Unix())
	for i := range nums {
		nums[i] = rand.Intn(1000)
	}

	var wg sync.WaitGroup
	wg.Add(DataSize / DataPerTask)
	tasks := make([]*Task, 0, DataSize/DataPerTask)
	for i := 0; i < DataSize/DataPerTask; i++ {
		task := &Task{
			index: i + 1,
			nums:  nums[i*DataPerTask : (i+1)*DataPerTask],
			wg:    &wg,
		}

		tasks = append(tasks, task)
		// 2.调用函数
		p.Invoke(task)
	}

	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())

	//验证数据
	var sum int
	for _, task := range tasks {
		sum += task.sum
	}

	var expect int
	for _, num := range nums {
		expect += num
	}

	fmt.Printf("finish all tasks, result is %d expect:%d\n", sum, expect)
}
