# Go 中监控代码性能
## 两个包：
1. net/http/pprof
使用场景：在线服务（一直运行着的程序）
2. runtime/pprof
使用场景：工具型应用（比如说定制化的分析小工具、集成到公司监控系统）
这两个包都是可以监控代码性能的， 只不过net/http/pprof是通过http端口方式暴露出来的，内部封装的仍然是runtime/pprof。

## net/http/pprof 源码分析
```go
//  net/http/pprof/pprof.go
func init() {
	http.HandleFunc("/debug/pprof/", Index)
	http.HandleFunc("/debug/pprof/cmdline", Cmdline)
	http.HandleFunc("/debug/pprof/profile", Profile)
	http.HandleFunc("/debug/pprof/symbol", Symbol)
	http.HandleFunc("/debug/pprof/trace", Trace)
}
```

Profile 举例
```go
func Profile(w http.ResponseWriter, r *http.Request) {
    // ...
	if err := pprof.StartCPUProfile(w); err != nil {
		// StartCPUProfile failed, so no writes yet.
		serveError(w, http.StatusInternalServerError,
			fmt.Sprintf("Could not enable CPU profiling: %s", err))
		return
	}
	sleep(w, time.Duration(sec)*time.Second)
	pprof.StopCPUProfile()
}
```
这个函数也是调用runtime/pprof的StartCPUProfile(w)方法开始 CPU profiling，然后睡眠一段时间（这个时间就是采样间隔），最后调用pprof.StopCPUProfile()停止采用.

StartCPUProfile()方法传入的是http.ResponseWriter类型变量，所以采样结果直接写回到 HTTP 的客户端

## 介绍：
runtime/pprof中的程序来生成三种包含实时性数据的概要文件，分别是
1. CPU概要文件   
在默认情况下，Go语言的运行时系统会以100 Hz的的频率对CPU使用情况进行取样。
2. 内存概要文件   
内存概要文件用于保存在用户程序执行期间的内存使用情况。这里所说的内存使用情况，其实就是程序运行过程中堆内存的分配情况。
3. 程序阻塞概要文件   
程序阻塞概要文件用于保存用户程序中的Goroutine阻塞事件的记录。



# 第三方性能分析来分析代码包
runtime.pprof 提供基础的运行时分析的驱动，但是这套接口使用起来还不是太方便，例如：
1. 输出数据使用 io.Writer 接口，虽然扩展性很强，但是对于实际使用不够方便，不支持写入文件。
2. 默认配置项较为复杂。

很多第三方的包在系统包 runtime.pprof 的技术上进行便利性封装，让整个测试过程更为方便。这里使用 github.com/pkg/profile 包进行例子展示，
使用下面代码安装这个包
```go
go get github.com/pkg/profile
```





