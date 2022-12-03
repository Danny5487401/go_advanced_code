# stack



在Go应用程序运行时，每个goroutine都维护着一个自己的栈区，这个栈区只能自己使用不能被其他goroutine使用。
栈区的初始大小是2KB（比x86_64架构下线程的默认栈2M要小很多），在goroutine运行的时候栈区会按照需要增长和收缩，占用的内存最大限制的默认值在64位系统上是1GB。
栈大小的初始值和上限这部分的设置都可以在Go的源码runtime/stack.go里找到：




## 参考资料
1. [解密Go协程的栈内存管理](https://mp.weixin.qq.com/s/ErnQDHeL5K8MPDYUPwjSYA)


