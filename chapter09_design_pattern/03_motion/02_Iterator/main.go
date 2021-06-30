package main

import "fmt"

/*
不使用迭代器的方案:
	首先要指出的是，绝大多数情况下Go程序是不需要用迭代器的。因为内置的slice和map两种容器都可以通过range进行遍历，并且这两种容器在性能方面做了足够的优化。
场景：
	当然某些特殊场合迭代器还是有用武之地。比如迭代器的Next()是个耗时操作，不能一口气拷贝所有元素；再比如某些条件下需要中断遍历
实现方式：
	推荐使用channel
标准库中递归实现：
	标准库中的container/ring中有Do()用法的例子
	type Ints []int

	func (i Ints) Do(fn func(int)) {
		for _, v := range i {
			fn(v)
		}
	}

	func main() {
		ints := Ints{1, 2, 3}
		ints.Do(func(v int) {
			fmt.Println(v)
		})
	}
*/

// channel实现
type Ints []int

func (i Ints) Iterator() <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range i {
			c <- v
		}
		close(c)
	}()
	return c
}
func main() {
	ints := Ints{1, 2, 3}
	for v := range ints.Iterator() {
		fmt.Println(v)
	}
}
