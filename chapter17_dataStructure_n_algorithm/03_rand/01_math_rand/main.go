package main

import (
	"fmt"
	"sync"

	"math/rand"
	"time"
)

func main() {

	// 生成相同的结果
	sameSourceRand()

	// 生成不同的结果
	newSourceRand()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Nanosecond * 10)
			fmt.Println(rrRand.Perm(10))
			wg.Done()
		}()
	}
	wg.Wait()

}

var (
	// 自定义生成Rand结构体，设置随机数种子
	rrRand = rand.New(rand.NewSource(time.Now().Unix()))
)

func sameSourceRand() {
	// 相同种子，每次运行的结果都是一样的
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().Unix())
		fmt.Printf("%v: 每次相同的结果：%v\n", time.Now().Unix(), rand.Intn(10))
	}
}

func newSourceRand() {
	// 相同种子，每次运行的结果都是一样的
	rand.Seed(time.Now().Unix())
	for i := 0; i < 5; i++ {
		fmt.Printf("%v: 每次不同的结果：%v\n", time.Now().Unix(), rand.Intn(10))
	}
}
