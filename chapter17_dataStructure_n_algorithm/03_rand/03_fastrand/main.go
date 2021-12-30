package main

import (
	"fmt"
	"github.com/NebulousLabs/fastrand"
)

// 伪随机的更优写法

func main() {
	//返回随机数(x~maxN)
	//fastrand.Uint64n(maxN uint32) uint64
	fmt.Println(fastrand.Uint64n(100))

}
