#happened before

简单来说就是如果事件 a 和事件 b 存在 happened-before 关系，即 a -> b，那么 a，b 完成后的结果一定要体现这种关系。
由于现代编译器、CPU 会做各种优化，包括编译器重排、内存重排等等，在并发代码里，happened-before 限制就非常重要了

根据晃岳攀老师在 Gopher China 2019 上的并发编程分享，关于 channel 的发送（send）、发送完成（send finished）、
接收（receive）、接收完成（receive finished）的 happened-before 关系如下：

    1.第 n 个 send 一定 happened before 第 n 个 receive finished，无论是缓冲型还是非缓冲型的 channel。
```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "hello world"
    // 发送
    done <- true
}


func main() {
    go aGoroutine()
    // 接受完成
	<-done
    fmt.Println(msg)
}

```
    
    2.对于容量为 m 的缓冲型 channel，第 n 个 receive 一定 happened before 第 n+m 个 send finished。

    缓冲型的 channel，当第 n+m 个 send 发生后，有下面两种情况:
        a.若第 n 个 receive 没发生。这时，channel 被填满了，send 就会被阻塞。那当第 n 个 receive 发生时，sender goroutine 会被唤醒，之后再继续发送过程。
        这样，第 n 个 receive 一定 happened before 第 n+m 个 send finished。
        
        b.若第 n 个 receive 已经发生过了，这直接就符合了要求。
    
    3.对于非缓冲型的 channel，第 n 个 receive 一定 happened before 第 n 个 send finished。
```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "hello world"
    // 接收
	<-done

}


func main() {
    go aGoroutine()
	// 发送完成
	done <- true

    fmt.Println(msg)
}

```
    
    4.channel close 一定 happened before receiver 得到通知
        回忆一下源码，先设置完 closed = 1，再唤醒等待的 receiver，并将零值拷贝给 receiver