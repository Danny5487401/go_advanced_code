package func_package

// .s 和 .go 文件的全局变量互通
var a = 999

func Get() int // 获取a的值

// Swap 简单返回值
func Swap(a, b int) (ret0, ret1 int)

// Foo 复杂返回值
func Foo(a bool, b int16) (c []byte) // 汇编函数声明
