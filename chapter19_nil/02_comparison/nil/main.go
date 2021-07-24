package main

import "fmt"

/*
理解为两个预定义标识符在进行比较，会报错，报错信息如下。但如果对nil进行重定义，标准的nil就会被覆盖，更改后的nil可以进行比较。
*/

func main() {
	//var nil = "不建议更改nil"  //可以比较
	fmt.Println(nil == nil)
}

// 报错
// invalid operation: nil == nil (operator == not defined on nil)
