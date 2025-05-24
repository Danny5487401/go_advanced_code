package main

import (
	"runtime"
	"syscall"
)

type File struct{ d int }

func main() {
	p := openFile("chapter06_pointer/07_keepalive_n_finalizer/01_wrong_use/t.txt")
	content := readFile(p.d) // 这里传入的是 int
	runtime.KeepAlive(p)

	println("Here is the content: " + content)
}

func openFile(path string) *File {
	d, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	p := &File{d}
	runtime.SetFinalizer(p, func(p *File) {
		syscall.Close(p.d) // 当变量被回收时，执行一些回收操作，加速一些资源的释放。在做性能优化的时候这样做确实有一定的效果，不过这样做是有一定的风险的。
	})

	return p
}

func readFile(descriptor int) string {
	doSomeAllocation()

	var buf [1000]byte
	_, err := syscall.Read(descriptor, buf[:])
	if err != nil {
		panic(err)
	}

	return string(buf[:])
}

func doSomeAllocation() {
	var a *int

	// memory increase to force the GC
	for i := 0; i < 10000000; i++ {
		i := 1
		a = &i
	}

	_ = a
}

/*
报错

panic: no such file or directory

goroutine 1 [running]:
main.openFile({0x10022d7fa?, 0x1001f6850?})
	/Users/python/Downloads/git_download/go_advanced_code/chapter06_pointer/07_keepalive/01_wrong_use/main.go:20 +0x84
main.main()
	/Users/python/Downloads/git_download/go_advanced_code/chapter06_pointer/07_keepalive/01_wrong_use/main.go:11 +0x28


*/
