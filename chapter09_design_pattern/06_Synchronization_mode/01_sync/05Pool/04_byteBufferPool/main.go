package main

/*
字节缓冲（底层一般是字节切片）
	在做字符串拼接时，为了拼接的高效，我们通常将中间结果存放在一个字节缓冲。在拼接完成之后，再从字节缓冲中生成结果字符串。在收发网络包时，也需要将不完整的包暂时存放在字节缓冲中。
标准库
	Go 标准库中的类型bytes.Buffer封装字节切片，提供一些使用接口。我们知道切片的容量是有限的，容量不足时需要进行扩容。而频繁的扩容容易造成性能抖动。
第三方库
	bytebufferpool实现了自己的Buffer类型，并使用一个简单的算法降低扩容带来的性能损失
*/
import (
	"fmt"

	"github.com/valyala/bytebufferpool"
)

func main() {
	// 使用默认的对象池
	b := bytebufferpool.Get()
	b.WriteString("hello")
	b.WriteByte(',')
	b.WriteString(" world!")

	fmt.Println(b.String())

	bytebufferpool.Put(b)
}

/*
细节：
	容量最小值取 2^6 = 64，因为这就是 64 位计算机上 CPU 缓存行的大小。这个大小的数据可以一次性被加载到 CPU 缓存行中，再小就无意义了。
	代码中多次使用atomic原子操作，避免加锁导致性能损失
缺点：
	由于大部分使用的容量都小于defaultSize，会有部分内存浪费
*/
