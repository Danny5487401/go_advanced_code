# 数值类型

标准库中的数值类型很多:

1. int/int8/int16/int32/int64
2. uint/uint8/uint16/uint32/uint64
3. float32/float64
4. byte/rune
   1. byte 是 uint8的别名,一个byte长度为8，即八位一个字节,一个byte等于八个bit,一个bit表示一位,8bit == byte
   2. alias for int32
5. uintptr
这些类型在汇编中就是一段存储着数据的连续内存，只是内存长度不一样，操作的时候看好数据长度就行。

   

## struct
struct 在汇编层面实际上就是一段连续内存，在作为参数传给函数时，会将其展开在 caller 的栈上传给对应的 callee:




