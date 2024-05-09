package main

import (
	"fmt"
	"github.com/oklog/run"
	"time"
)

// 定义了三个 actor，前两个 actor 一直等待。第三个 actor 在 3s 后结束退出。引起前两个 actor 退出。
func main() {
	g := run.Group{}
	{
		cancel := make(chan struct{})
		g.Add(
			func() error {

				select {
				case <-cancel:
					fmt.Println("Go routine 1 is closed")
					break
				}

				return nil
			},
			func(error) {
				close(cancel)
			},
		)
	}
	{
		cancel := make(chan struct{})
		g.Add(
			func() error {

				select {
				case <-cancel:
					fmt.Println("Go routine 2 is closed")
					break
				}

				return nil
			},
			func(error) {
				close(cancel)
			},
		)
	}
	{
		g.Add(
			func() error {
				for i := 0; i < 3; i++ {
					time.Sleep(1 * time.Second)
					fmt.Println("Go routine 3 is sleeping...")
				}
				fmt.Println("Go routine 3 is closed")
				return nil
			},
			func(error) {
				return
			},
		)
	}
	g.Run()
}
