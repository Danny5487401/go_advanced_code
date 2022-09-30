package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, request *http.Request) {
		// 假设异常

		// 第一种情况
		//panic("出现异常")
		//a:=10
		//b:=0
		//fmt.Println(a/b) //出现异常，页面没有打印出来，但是程序并没有挂掉，多线程

		// 第二种情况： 线程里面 嵌套 线程：注意
		go func() {
			// 例如：操作数据库如redis,有人觉得这段代码可以放在协程中，有很大隐患

			// 解决方法：
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("捕获到异常:%v\n", err)
				}
			}()

			// panic 会引起主线程的挂掉，同时导致其他的协程也挂了
			panic("出现异常") //页面打印出来，但是程序挂了

			// 原因：父协程无法捕获子协程中出现的异常

		}()

		w.Write([]byte("hello world"))
	})
	http.ListenAndServe("127.0.0.1:9090", nil) // 内部注册了recover
}
