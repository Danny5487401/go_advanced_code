package main

/*
第三方性能分析来分析代码包
	runtime.pprof 提供基础的运行时分析的驱动，但是这套接口使用起来还不是太方便，例如：
		1.输出数据使用 io.Writer 接口，虽然扩展性很强，但是对于实际使用不够方便，不支持写入文件。
		2.默认配置项较为复杂。

很多第三方的包在系统包 runtime.pprof 的技术上进行便利性封装，让整个测试过程更为方便。这里使用 github.com/pkg/profile 包进行例子展示，
使用下面代码安装这个包
	go get github.com/pkg/profile
*/
