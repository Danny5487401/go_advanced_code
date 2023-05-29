<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [compile 编译](#compile-%E7%BC%96%E8%AF%91)
  - [选项](#%E9%80%89%E9%A1%B9)
  - [其他的编译器指令](#%E5%85%B6%E4%BB%96%E7%9A%84%E7%BC%96%E8%AF%91%E5%99%A8%E6%8C%87%E4%BB%A4)
    - [//go:noescape](#gonoescape)
    - [//go:uintptrescapes](#gouintptrescapes)
    - [//go:noinline](#gonoinline)
    - [//go:norace](#gonorace)
    - [//go:nosplit](#gonosplit)
    - [//go:linkname localname [importpath.name]](#golinkname-localname-importpathname)
  - [参考链接](#%E5%8F%82%E8%80%83%E9%93%BE%E6%8E%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# compile 编译 

compile（通常作为“go tool compile”调用）编译多个文件组成的Go包。然后，写入一个对象文件，以第一个源码文件名的.o的的中间目标文件 (intermediate object file)。
然后，可以将对象文件与其他对象组合到包归档文件中，或者直接传递到链接器（“go tool link”）。 如果使用-pack调用，编译器将跳过中间对象文件，直接写入归档文件 (package archive)。

生成的目标文件包含包本身暴露的符号的类型信息， 也包含包引用的其他包的符号的类型信息。 所以在编译调用一个包的包时， 只需要读取被调用的包的目标文件即可。

```shell
go tool compile [flags] file...

```
file一定需要是一整个package所有的代码文件。
## 选项

```shell
➜  03_n git:(feature/memory) ✗ go tool compile --help    
usage: compile [options] file.go...
  -% int
        debug non-static initializers
  -+    compiling runtime
  -B    disable bounds checking
  -C    disable printing of columns in error messages
  -D path 用于本地引用依赖的相对路径
  -E    debug symbol export
  -G    accept generic code (default 3)
  -I directory 在查询完$GOROOT/pkg/$GOOS_$GOARCH之后， 进一步从dir1和dir2查询需要的依赖包。
  -K    debug missing line numbers
  -L    在错误信息中展示完整的文件路径
  -N    禁用优化
  -S    将code的汇编输出到标准输出
  -V     输出编译器的版本
  -W    debug parse tree after type checking
  -asan
        build code compatible with C/C++ address sanitizer
  -asmhdr file 将汇编的头写到file中
  -bench file
        append benchmark times to file
  -blockprofile file
        write block profile to file
  -buildid id  将id作为build_id写入输出的元数据中
  -c int  编译时并发度，默认是1，表示不并发
  -clobberdead
        clobber dead stack slots (for debugging)
  -clobberdeadreg
        clobber dead registers (for debugging)
  -complete 假定包不含有非go的部分
  -cpuprofile file 将编译时的CPU profile写入到file中
  -d value
        enable debugging settings; try -d help
  -dwarf
        generate DWARF symbols (default true)  生成DWARF符号
  -dwarfbasentries
        use base address selection entries in DWARF
  -dwarflocationlists
        add location lists to DWARF in optimized mode (default true) 优化模式中， 向DWARF增加位置列表 （location list）
  -dynlink 允许引用在共享库中的go符号
  -e    no limit on number of errors reported 移除错误数量的上限
  -embedcfg file
        read go:embed configuration from file
  -gendwarfinl int
        generate DWARF inline info records (default 2)
  -goversion string 定义使用的go tool版本， 用于runtime的版本和goversion不匹配的情况
  -h    halt on error 当第一个错误被发现时停止，并输出堆栈trace
  -importcfg file 从file读取配置。 配置包含importmap/packagefile
  -importmap definition
        add definition of the form source=actual to import map 在编译时，将对old的引用更换为new。 这个flag可以有多个来设置多个映射
  -installsuffix suffix
        set pkg directory suffix 从$GOROOT/pkg/$GOOS_$GOARCH_suffix查找包， 而不是$GOROOT/pkg/$GOOS_$GOARCH
  -j    debug runtime-initialized variables
  -json string
        version,file for JSON compiler/optimizer detail output
  -l    禁用内联
  -lang string 使用的语言版本，如-lang=go1.12， 默认使用当前版本
  -linkobj file
        write linker-specific object to file
  -linkshared
        generate code that will be linked against Go shared libraries
  -live
        debug liveness analysis
  -m    输出优化决定。可以传入整数 (-m=10) ，越大的整数输出越详细
  -memprofile file 输出编译时的内存profile到file
  -memprofilerate rate
        set runtime.MemProfileRate to rate 设置编译时的runtime.MemProfileRate为rate
  -msan
        build code compatible with C/C++ memory sanitizer 开启内存检查器 （memory sanitizer）
  -mutexprofile file
        write mutex profile to file  编译时的mutex profile写入到file
  -nolocalimports 禁用本地引用/相对引用。
  -o file 将目标文件写入到file
  -p path
        set expected package import path 判断如果增加了对于path的引用是否会出现循环引用的问题
  -pack
        write to file.a instead of file.o 输出打包过的格式，而不是目标文件
  -r    debug generated wrappers
  -race 开启数据竞争检测
  -shared
        generate code that can be linked into a shared library 生成可以链接到共享库的代码
  -smallframes
        reduce the size limit for stack allocated objects
  -spectre list
        enable spectre mitigations in list (all, index, ret)
  -std
        compiling standard library
  -symabis file
        read symbol ABIs from file
  -t    enable tracing for debugging the compiler
  -traceprofile file
        write an execution trace to file
  -trimpath prefix
        remove prefix from recorded source file paths
  -v    increase debug verbosity
  -w    debug type checking
  -wb
        enable write barrier (default true)

```

## 其他的编译器指令

其他的编译器指令都是由//go:name的形式。


### //go:noescape

noescape后续紧跟着一个没有body的函数声明 （没有body意味着这个函数是用非go实现的）。 noescape意味着这个函数不允许接受任何会逃逸到堆上的指针作为参数， 或者逃逸到该函数的返回值中
```go
// /Users/python/go/pkg/mod/github.com/cespare/xxhash/v2@v2.1.2/xxhash_amd64.go
// +build !appengine
// +build gc
// +build !purego

package xxhash

// Sum64 computes the 64-bit xxHash digest of b.
//
//go:noescape
func Sum64(b []byte) uint64
```

### //go:uintptrescapes
uintptrescapes后续需要紧跟一个函数声明。 指令表明该函数的uintptr可能是由指向被垃圾回收管理的对象的指针转换而来的


```go
// /Users/python/go/go1.18/src/syscall/syscall_linux.go

// AllThreadsSyscall performs a syscall on each OS thread of the Go
// runtime. It first invokes the syscall on one thread. Should that
// invocation fail, it returns immediately with the error status.
// Otherwise, it invokes the syscall on all of the remaining threads
// in parallel. It will terminate the program if it observes any
// invoked syscall's return value differs from that of the first
// invocation.
//
// AllThreadsSyscall is intended for emulating simultaneous
// process-wide state changes that require consistently modifying
// per-thread state of the Go runtime.
//
// AllThreadsSyscall is unaware of any threads that are launched
// explicitly by cgo linked code, so the function always returns
// ENOTSUP in binaries that use cgo.
//go:uintptrescapes
func AllThreadsSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	if cgo_libc_setegid != nil {
		return minus1, minus1, ENOTSUP
	}
	r1, r2, errno := runtime_doAllThreadsSyscall(trap, a1, a2, a3, 0, 0, 0)
	return r1, r2, Errno(errno)
}
```

### //go:noinline
表明这个函数不可以被内联

### //go:norace

这个函数的数据访问应该被数据竞争探测器忽略。

### //go:nosplit
表明这个函数需要忽略通常的堆栈溢出检查。 用于这个函数被调用时，调用者goroutine不可以被抢占

### //go:linkname localname [importpath.name]
编译器在目标文件的符号使用importpath.name替换掉localname。 如果[importpath.name]被忽略，那么会使用默认的符号名， 但产生副作用：使localname可以被其他包访问

```shell
// /Users/python/go/go1.18/src/time/time.go

// runtimeNano returns the current value of the runtime clock in nanoseconds.
//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64
```


## 参考链接

1. [官方compile](https://pkg.go.dev/cmd/compile#hdr-Command_Line)