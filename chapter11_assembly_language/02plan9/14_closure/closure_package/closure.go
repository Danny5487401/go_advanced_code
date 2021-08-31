package closure_package

import "unsafe"

/*
// NewTwiceFunClosure 闭包函数对象捕获了外层的x参数。返回的闭包函数对象在执行时，每次将捕获的外层变量乘以2之后再返回
func NewTwiceFunClosure(x int) func() int {
	return func() int {
		x *= 2
		return x
	}
}

func main() {
	fnTwice := NewTwiceFunClosure(1)

	println(fnTwice()) // 1*2 => 2
	println(fnTwice()) // 2*2 => 4
	println(fnTwice()) // 4*2 => 8
}


*/

func ptrToFunc(p unsafe.Pointer) func() int

// 用于返回闭包函数机器指令的开始地址（类似全局函数的地址）
func asmFunTwiceClosureAddr() uintptr

// 是闭包函数对应的全局函数的实现
func asmFunTwiceClosureBody() int

// 手动构造闭包
type FunTwiceClosure struct {
	F uintptr // 表示闭包函数的函数指令的地址
	X int     // 闭包捕获的外部变量
}

func NewTwiceFunClosure(x int) func() int {
	var p = &FunTwiceClosure{
		F: asmFunTwiceClosureAddr(), // asmFunTwiceClosureAddr函数用于辅助获取闭包函数的函数指令的地址，采用汇编语言实现
		X: x,
	}
	// 将结构体指针转为闭包函数对象返回，该函数也是通过汇编语言实现
	return ptrToFunc(unsafe.Pointer(p))
}
