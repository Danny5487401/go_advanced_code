package func_package

var a = 999

func Get() int

// Swap 简单返回值
func Swap(a, b int) (ret0, ret1 int)

// Foo 复杂返回值
func Foo(a bool, b int16) (c []byte)

// 调用函数时，被调用函数的参数和返回值内存空间都必须由调用者提供。
func printsum(a, b int) {
	var ret = sum(a, b)
	println(ret)
}

func sum(a, b int) int {
	return a + b
}
