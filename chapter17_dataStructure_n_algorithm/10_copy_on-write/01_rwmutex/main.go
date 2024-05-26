package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 读写锁
var mu sync.RWMutex

type CowMap map[int]string

func (c *CowMap) Set(key int, value string) {
	(*c)[key] = value
}

func (c *CowMap) Get(key int) string {
	return (*c)[key]
}

func readLoop(c *CowMap) {
	for {
		//读的时候上读锁
		mu.RLock()
		fmt.Println(c.Get(3))
		mu.RUnlock()
	}
}

func writeLoop(c *CowMap) {
	for i := 0; i < 10000000; i++ {
		//每隔5s写一次
		time.Sleep(time.Second * 5)
		//写的时候上写锁
		mu.Lock()
		c.Set(3, "werben-"+strconv.Itoa(i))
		mu.Unlock()
	}
}

func main() {
	c := make(CowMap)
	c.Set(1, "a")
	c.Set(2, "b")
	c.Set(3, "c")

	go readLoop(&c)
	writeLoop(&c)
}

/*
只是每隔5秒写一次，但是上面的读锁还是一直不断的上锁解锁，这个在没有写数据的时候，其实都是没有意义的。如果时间更长，比如1天才修改一次，读锁浪费了大量的无用资源
*/
