package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
)

type tmap map[int]string

type CowMap struct {
	mu   sync.Mutex
	data atomic.Value
}

func NewCowMap() *CowMap {
	m := make(tmap)
	c := &CowMap{}
	c.data.Store(m)
	return c
}

func (c *CowMap) clone() tmap {
	m := make(tmap)
	for k, v := range c.data.Load().(tmap) {
		m[k] = v
	}
	return m
}

func (c *CowMap) Get(key int) (value string, ok bool) {
	value, ok = c.data.Load().(tmap)[key]
	return
}

func (c *CowMap) Set(key int, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//写之前，先拷贝一个副本
	copyOn := c.clone()
	//修改副本
	copyOn[key] = value
	//修改副本数据后，将副本转正
	c.data.Store(copyOn)
}

func readLoop(c *CowMap) {
	for {
		fmt.Println(c.Get(3))
	}
}

func writeLoop(c *CowMap) {
	for i := 0; i < 10000000; i++ {
		c.Set(3, "werben-"+strconv.Itoa(i))
	}
}

func main() {
	c := NewCowMap()
	c.Set(1, "a")
	c.Set(2, "b")
	c.Set(3, "c")

	go readLoop(c)
	writeLoop(c)
}
