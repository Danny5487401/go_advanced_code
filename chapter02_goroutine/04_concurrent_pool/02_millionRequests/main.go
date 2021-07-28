package main

import (
	"fmt"
	"runtime"
	"time"
)

//Job 任务
type Job interface {
	Do() //do something
}

// worker 工人
type Worker struct {
	JobQueue chan Job  //任务队列
	Quit     chan bool //停止当前任务
}

// 新建一个worker实例
func NewWorker() Worker {
	return Worker{
		JobQueue: make(chan Job), //初始化工作队列为null
		Quit:     make(chan bool),
	}
}

/*
整个过程中 每个Worker(工人)都会被运行在一个协程中，
在整个WorkerPool(领导)中就会有num个可空闲的Worker(工人)，
当来一条数据的时候，领导就会小组中取一个空闲的Worker(工人)去执行该Job，
当工作池中没有可用的worker(工人)时，就会阻塞等待一个空闲的worker(工人)。
每读到一个通道参数 运行一个 worker
*/

func (w Worker) Run(wq chan chan Job) {
	//这是一个独立的协程 循环读取通道内的数据，
	//保证 每读到一个通道参数就 去做这件事，没读到就阻塞
	go func() {
		for {
			wq <- w.JobQueue //注册工作通道  到 线程池
			select {
			case job := <-w.JobQueue: //读到参数
				job.Do()
			case <-w.Quit:
				return
			}
		}
	}()
}

// workerPool 领导
type WorkerPool struct {
	WorkerLen   int           //线程池中 worker工人数量
	JobQueue    chan Job      //线程池中Job通道
	WorkerQueue chan chan Job // 通道里放chan Job
}

func NewWorkerPool(workerLen int) *WorkerPool {
	return &WorkerPool{
		WorkerLen:   workerLen,                      //开始建立 workerLen 个worker(工人)协程
		JobQueue:    make(chan Job),                 //工作队列 通道
		WorkerQueue: make(chan chan Job, workerLen), //最大通道参数设为 最大协程数 workerLen 工人的数量最大值
	}
}

func (wp WorkerPool) Run() {
	//初始化时会按照传入的num，启动num个后台协程，然后循环读取Job通道里面的数据，
	//读到一个数据时，再获取一个可用的Worker，并将Job对象传递到该Worker的chan通道
	fmt.Println("初始化Worker")
	for i := 0; i < wp.WorkerLen; i++ {
		//新建 workerLen  个 worker(工人) 协程(并发执行)，每个协程可处理一个请求
		worker := NewWorker() //运行一个协程 将线程池 通道的参数  传递到 worker协程的通道中 进而处理这个请求
		worker.Run(wp.WorkerQueue)
	}

	// 循环获取可用的worker,往worker中写job
	go func() { //这是一个单独的协程 只负责保证 不断获取可用的worker
		for {
			select {
			case job := <-wp.JobQueue: //读取任务
				//尝试获取一个可用的worker作业通道。
				//这将阻塞，直到一个worker空闲
				worker := <-wp.WorkerQueue
				worker <- job //将任务 分配给该工人
			}
		}
	}()

}

// 模拟请求
type DoSomething struct {
	Num int
}

// 保证任务有Do方法
func (d *DoSomething) Do() {
	fmt.Println("开启线程数：", d.Num)
	//这里延迟是模拟处理数据的耗时
	time.Sleep(1 * 1 * time.Second)
}

func main() {
	//设置最大线程数
	num := 100 * 100 * 20

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)并发个数 20万个
	p := NewWorkerPool(num)
	p.Run() //有任务就去做，没有就阻塞，任务做不过来也阻塞

	//dataNum := 100 * 100 * 100    //模拟百万请求
	dataNum := 100 * 100 * 100
	go func() {
		//这是一个独立的协程 保证可以接受到每个用户的请求
		for i := 0; i <= dataNum; i++ {
			sc := &DoSomething{i}
			p.JobQueue <- sc //往线程池 的通道中 写参数   每个参数相当于一个请求  来了100万个请求
		}

	}()
	//这里循环打印查看协程个数
	for { //阻塞主程序结束
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}
}
