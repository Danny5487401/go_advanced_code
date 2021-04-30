package _3Golang_data_structure

/*
数值类型

标准库中的数值类型很多:

1. int/int8/int16/int32/int64
2. uint/uint8/uint16/uint32/uint64
3. float32/float64
4. byte/rune
5. uintptr
这些类型在汇编中就是一段存储着数据的连续内存，只是内存长度不一样，操作的时候看好数据长度就行。

slice

slice 在传递给函数的时候，实际上会展开成三个参数:
1. 首元素地址
2. slice 的 len
3. slice 的 cap

string_test

//go:noinline
func stringParam(s string_test) {}

func main() {
    var x = "abcc"
    stringParam(x)
}
go tool compile -S 输出其汇编:
在汇编层面 string_test 就是地址 + 字符串长度。


struct
struct 在汇编层面实际上就是一段连续内存，在作为参数传给函数时，会将其展开在 caller 的栈上传给对应的 callee:
 */



