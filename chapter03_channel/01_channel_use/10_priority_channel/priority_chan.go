package main

import "fmt"

/*
现在我们有一个需求：我们有一个函数会持续不间断地从ch1和ch2中分别接收任务1和任务2，

如何确保当ch1和ch2同时达到就绪状态时，优先执行任务1，在没有任务1的时候再去执行任务2呢？
*/

// 方式一：
func worker1(ch1, ch2 <-chan int, stopCh chan struct{}) {

	for {
		select {
		case <-stopCh:
			return
		case job1 := <-ch1:
			fmt.Println(job1)
		default:
			select {
			case job2 := <-ch2:
				fmt.Println(job2)
			default:
			}
		}
	}
}

// 方式一缺点：worker1,如果ch1和ch2都没有达到就绪状态的话，整个程序不会阻塞而是进入了死循环

// 方式二
func worker2(ch1, ch2 <-chan int, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		case job1 := <-ch1:
			fmt.Println(job1)
		case job2 := <-ch2:
		priority:
			for {
				select {
				case job1 := <-ch1:
					fmt.Println(job1)
				default:
					break priority
				}
			}
			fmt.Println(job2)
		}
	}
}
