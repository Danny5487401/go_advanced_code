package main

import (
	"fmt"
)

func main() {

	s := []int{10, 20, 30}
	fmt.Printf("%p,%+v\n", s, s)

	changeSlice(s)

	fmt.Printf("%p,%+v\n", s, s)
}

func changeSlice(s []int) {

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	fmt.Printf("%p,%+v\n", s, s)
}

// 运行结果
//0x14000124000,[10 20 30]
//0x1400012a000,[10 20 30 0 1 2 3 4 5 6 7 8 9]
//0x14000124000,[10 20 30]
