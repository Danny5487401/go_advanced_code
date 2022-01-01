package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 优雅退出go守护进程
func main() {
	// 创建监听退出chan
	c := make(chan os.Signal)
	// 监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				ExitFunc()
			case syscall.SIGUSR1:
				fmt.Println("usr1", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2", s)
			default:
				fmt.Println("other", s)
			}
		}
	}()

	fmt.Println("进程启动...")
	sum := 0
	for {
		sum++
		fmt.Println("sum:", sum)
		time.Sleep(time.Second)
	}
}

func ExitFunc() {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	fmt.Println("结束退出...")
	os.Exit(0)
}

/*

	ctrl+c退出,输出
	退出信号 interrupt

	kill pid 输出
	退出信号 terminated

	kill -USR1 pid 输出
	退出信号 user defined signal 1

	kill -USR2 pid 输出
	退出信号 user defined signal 2
*/
