package main

import (
	"fmt"
	"sync"
)

/*
单例模式：
	始化单例资源，或者并发访问只需初始化一次的共享资源，或者在测试的时候初始化一次测试资源。
分析：
	sync.Once 只暴露了一个方法 Do，你可以多次调用 Do 方法，但是只有第一次调用 Do 方法时 f 参数才会执行，这里的 f 是一个无参数无返回值的函数。
*/

func main() {
	// 实例1

	//// 定义once结构体对象
	//var once sync.Once
	//
	//// 需要只执行一次的函数
	//onceBody := func() {
	//	fmt.Println("Only once")
	//}
	//
	//done := make(chan bool)
	//for i := 0; i < 10; i++ {
	//	fmt.Println(i)
	//	go func() {
	//		once.Do(onceBody)
	//		done <- true
	//	}()
	//}
	//
	//for i := 0; i < 10; i++ {
	//	<-done
	//}

	// 实例2

	// 如果多次调用once.Do(f)，则f仅执行一次，即使f在每次调用中为不同的值。如果你想多次执行的话就需要多个sync.Once
	var once1 sync.Once
	onceBody1 := func() {
		fmt.Println("only one1")
	}

	onceBody2 := func() {
		fmt.Println("only one2")
	}

	done1 := make(chan bool)
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		go func() {
			once1.Do(onceBody1)
			once1.Do(onceBody2)
			done1 <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done1
	}

}

/* 结果
0
1
2
3
Only once
4
5
6
7
8
9

*/

/*
源码分析:// sync/once.go

type Once struct {
   done uint32 // 初始值为0表示还未执行过，1表示已经执行过
   m    Mutex
}
func (o *Once) Do(f func()) {
   // 判断done是否为0，若为0，表示未执行过，调用doSlow()方法初始化
   if atomic.LoadUint32(&o.done) == 0 {
      // Outlined slow-path to allow inlining of the fast-path.
      o.doSlow(f)
   }
}

// 加载资源
func (o *Once) doSlow(f func()) {
   o.m.Lock()
   defer o.m.Unlock()
   // 采用双重检测机制 加锁判断done是否为零
   if o.done == 0 {
      // 执行完f()函数后，将done值设置为1
      defer atomic.StoreUint32(&o.done, 1)
      // 执行传入的f()函数
      f()
   }
}

*/
