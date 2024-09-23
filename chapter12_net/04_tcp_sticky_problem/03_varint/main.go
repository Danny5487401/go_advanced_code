package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	buf := make([]byte, binary.MaxVarintLen64)
	for _, x := range []int64{-65, 1, 2, 127, 128, 255, 256} {
		n := binary.PutVarint(buf, x)
		fmt.Printf("%v 输出的可变长度为：%v 十六进制为：%x\n", x, n, buf[:n])
	}
}
