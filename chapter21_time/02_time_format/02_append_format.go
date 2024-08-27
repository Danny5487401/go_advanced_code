package main

import (
	"fmt"
	"time"
)

func main() {
	var a [64]byte
	var b = a[:0]                                      // 使用a[:0]的表示法将固定大小的数组转换为由此数组所支持的b表示的切片类型。这样可以通过编译器检查，并且会在栈上面分配内存。
	byteInfo := time.Now().AppendFormat(b, "20060102") // AppendFormat()，这个方法本身通过编译器栈分配检查。而之前版本Format()，编译器不能确定需要分配的内存大小，所以不满足栈上分配规则
	fmt.Println(string(byteInfo))

}
