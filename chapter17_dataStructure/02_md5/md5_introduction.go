package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

/*
md5
	md5算法属于hash算法的一种，所以在了解md5之前，我们先认识一下go提供的hash接口。hash算法是保证只要输入的值不同，就一定会得到两个不同的指定长度的hash值。当前两个不同值产生相同的hash还是有可能的，只是这个可能性很小很小
hash源码:hash包
	type Hash interface {
		// 通过io.Writer接口的Write方法向hash中添加数据
		io.Writer
		// 返回添加b到当前的hash值后的新切片，不会改变底层的hash状态，这个方法就是返回计算后的hash值，只是它是字符切片
		Sum(b []byte) []byte
		// 重设hash为无数据输入的状态，就是清空hash之前写入的数据
		Reset()
		// 返回Sum会返回的切片的长度
		Size() int
		// 返回hash底层的块大小；Write方法可以接受任何大小的数据，
		// 但提供的数据是块大小的倍数时效率更高
		BlockSize() int
	}
 Hash包还有两个Hash接口：
type Hash32 interface { // Hash32是一个被所有32位hash函数实现的公共接口。
    Hash
    Sum32() uint32
}
type Hash64 interface { // Hash64是一个被所有64位hash函数实现的公共接口。
    Hash
    Sum64() uint64
}
md5实现的Hash接口是16位的hash函数（它的Sum方法返回的字符切片长度为16位），Hash32和hash64是属于安全性更高的两个Hash函数，产生的hash值也更长。
*/

func main() {
	//// 方式一
	// 1. 创建md5算法
	has := md5.New()
	// 2. 写入需要加密的数据
	has.Write([]byte("abc123"))
	b := has.Sum(nil) // 获取hash值字符切片；Sum函数接受一个字符切片，这个切片的内容会原样的追加到abc123加密后的hash值的前面，这里我们不需要这么做，所以传入nil
	fmt.Println(b)    // 打印一下 [233 154 24 196 40 203 56 213 242 96 133 54 120 146 46 3]
	// 上面可以看到加密后的数据为长度为16位的字符切片，一般我们会把它转为16进制，方便存储和传播，下一步转换16进制
	fmt.Println(hex.EncodeToString(b)) // 通过hex包的EncodeToString函数，将数据转为16进制字符串； 打印 e99a18c428cb38d5f260853678922e03

	// 还有一种方法转换为16进制,通过fmt的格式化打印方法， %x表示转换为16进制
	fmt.Printf("%x", b) // 打印 e99a18c428cb38d5f260853678922e03

	//方式二
	//func Sum(data []byte) [Size]byte直接返回数据data的MD5加密值，注意它返回的是指定大小(Size)的数组，而不是切片了
	//b := md5.Sum([]byte("abc123")) // 加密数据
	//fmt.Printf("%x", b)            // 转换为16进制，并打印
}
