package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func main() {
	// 方式一
	// 1. 创建md5算法
	h := md5.New()

	// 2. 写入需要加密的数据
	h.Write([]byte("abc123"))
	b := h.Sum(nil) // 获取hash值字符切片；Sum函数接受一个字符切片，这个切片的内容会原样的追加到abc123加密后的hash值的前面，这里我们不需要这么做，所以传入nil
	fmt.Println(b)  // 打印一下 [233 154 24 196 40 203 56 213 242 96 133 54 120 146 46 3]

	// 3. 打印结果
	// 上面可以看到加密后的数据为长度为16位的字符切片，一般我们会把它转为16进制，方便存储和传播，下一步转换16进制
	// a.通过hex包的EncodeToString函数，将数据转为16进制字符串； 打印 e99a18c428cb38d5f260853678922e03
	fmt.Println(hex.EncodeToString(b))
	// b.还有一种方法转换为16进制,通过fmt的格式化打印方法， %x表示转换为16进制
	fmt.Printf("%x", b) // 打印 e99a18c428cb38d5f260853678922e03----这样打印性能差

	// 方式二:简化写法
	//func Sum(data []byte) [Size]byte直接返回数据data的MD5加密值，注意它返回的是指定大小(Size)的数组，而不是切片了
	b2 := md5.Sum([]byte("abc123")) // 加密数据
	fmt.Printf("%x", b2)            // 转换为16进制，并打印

	// 方式三: 使用io库
	fmt.Println(Md5UsingIo("abc123"))
}

func Md5UsingIo(s string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, s)
	return hex.EncodeToString(hash.Sum(nil))
}
