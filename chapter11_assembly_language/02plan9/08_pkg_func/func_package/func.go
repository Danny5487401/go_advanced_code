package func_package

var a = 999

func Get() int

// Swap 简单返回值
func Swap(a, b int) (ret0, ret1 int)

// Foo 复杂返回值
func Foo(a bool, b int16) (c []byte)

func add(a, b int) int // 汇编函数声明

func sub(a, b int) int // 汇编函数声明

func mul(a, b int) int // 汇编函数声明
