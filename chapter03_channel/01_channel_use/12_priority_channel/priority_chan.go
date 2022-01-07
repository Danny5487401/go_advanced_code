package _2_priority_channel

import "fmt"

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

// Note:worker1,如果ch1和ch2都没有达到就绪状态的话，整个程序不会阻塞而是进入了死循环

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
