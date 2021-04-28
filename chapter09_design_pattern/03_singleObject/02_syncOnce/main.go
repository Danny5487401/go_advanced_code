package main

import (
	"fmt"
	"sync"
)

// 对于单例模式来说Golang在标准包中已经有了更好的解决方法Sync.Once



type singleton struct{}
var ins *singleton
var once sync.Once
func GetIns() *singleton {
	once.Do(func(){
		ins = &singleton{}
	})
	return ins
}


/* 源码
1. once 结构
type Once struct {
	m    Mutex
	done uint32
}

2. 主要方法

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

 */


// 案例

func main() {
	var once sync.Once
	onceBody1 := func() {
		fmt.Println("only one Body1")
	}

	onceBody2 := func() {
		fmt.Println("only one Body2")
	}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		go func() {
			once.Do(onceBody1)
			once.Do(onceBody2)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// 注意：次调用once.Do(f)，则f仅执行一次，即使f在每次调用中为不同的值。如果你想多次执行的话就需要多个sync.Once

