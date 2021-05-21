package main

// 滥用 defer 可能会导致性能问题，尤其是在一个 "大循环" 里。
import (
	"fmt"
	"sync"
	"time"
)

var lock sync.Mutex

func test() {
	lock.Lock()
	lock.Unlock()
}

func testDefer() {
	lock.Lock()
	defer lock.Unlock()
}

func main() {
	func() {
		t1 := time.Now()

		for i := 0; i < 100000; i++ {
			test()
		}
		elapsed := time.Since(t1)
		fmt.Println("test elapsed: ", elapsed)
	}()
	func() {
		t1 := time.Now()

		for i := 0; i < 100000; i++ {
			testDefer()
		}
		elapsed := time.Since(t1)
		fmt.Println("testDefer elapsed: ", elapsed)
	}()

}
