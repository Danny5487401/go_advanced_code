package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// 下面这些高阶函数都是使用的全局源
	fmt.Printf("%v\n", rand.Int63())             // 产生一个0到2的63次方之间的正整数,类型是int64
	fmt.Printf("每次相同的结果：%v\n", rand.Int63n(100)) // 产生一个0到n之间的正整数,n不能大于2的63次方，类型是int64
	fmt.Printf("%v\n", rand.Int31())             // 这个和前面的区别就产生的是int32类型的整数
	fmt.Printf("%v\n", rand.Int31n(100))         // 同上

	// 这个比较有意思了，它至少产生一个32位的正整数，因为int类型在64位机器上等于int64，在32位机器上就是int32
	fmt.Printf("%v\n", rand.Int())
	fmt.Printf("%v\n", rand.Intn(100)) // 同上

	fmt.Printf("%v\n", rand.Float32()) // 产生一个0到1.0的浮点数，float32类型
	fmt.Printf("%v\n", rand.Float64()) // 产生一个0到1.0的浮点数，float64类型

	// 生成不同的结果
	newSourceRand()

}

func newSourceRand() {
	//相同种子，每次运行的结果都是一样的
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		fmt.Printf("current:%d\n", time.Now().Unix())
		rand.Seed(time.Now().Unix())
		fmt.Printf("每次不同的结果：%v\n", rand.Intn(100))
	}
}
