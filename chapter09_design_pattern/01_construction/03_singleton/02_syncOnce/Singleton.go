package main

import (
	"fmt"
	"sync"
)

type Singleton struct{}

var (
	singleton *Singleton
	once      sync.Once
)

func GetSingleton() *Singleton {
	once.Do(func() {
		fmt.Println("Create Obj")
		singleton = new(Singleton)
	})
	return singleton
}
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			obj := GetSingleton()
			fmt.Printf("%p\n", obj) //打印其地址
			wg.Done()
		}()
	}
	wg.Wait()
	// 只打印一次 Create Obj
}
