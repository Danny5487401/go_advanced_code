package main

import (
	"fmt"
	"sync"
)

func main() {

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
