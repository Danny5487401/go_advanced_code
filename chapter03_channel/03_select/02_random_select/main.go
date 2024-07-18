package main

func main() {
	ch := make(chan int)
	done := make(chan struct{})
	go func() {
		for i := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
			ch <- i // 注意这里使用无缓存 chan，所以是阻塞的
		}
		close(done)
	}()

	for {
		select {
		case <-done:
			println("done")
			return
		case <-ch:
			println("case1")
		case <-ch:
			println("case2")

		}
	}
}
