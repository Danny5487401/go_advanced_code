package main

// #include <pthread.h>
import "C"

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	fmt.Println("main", C.pthread_self())
	go func() {
		runtime.LockOSThread()
		fmt.Println("locked", C.pthread_self())
		go func() {
			fmt.Println("locked child", C.pthread_self())
			ch1 <- true
		}()
		ch2 <- true
	}()
	<-ch1
	<-ch2
}

/* linux amd64
main 140380350363456
locked 140380350363456
locked child 140379674474240

*/
