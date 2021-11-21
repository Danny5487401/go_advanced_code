// io 包这些接口和原始的操作以不同的实现包装了低级操作，客户不应假定它们对于并行执行是安全的
// 在io包中最重要的是两个接口：Reader和Writer接口

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	/*
		读取数据：
			Reader接口：
				Read(p []byte)(n int, error)
	*/
	//1.os库用法：读取本地aa.txt文件中的数据
	osRead()

	// 2。bytes库用法
	byteRead()

}

/*
bytes.Buffer源码分析
	type Buffer struct {
		buf      []byte // contents are the bytes buf[off : len(buf)]   缓冲区
		off      int    // 写的索引值，指针偏移量
		lastRead readOp // 上次读的操作，用于UnReadRune等撤回到上次读操作之前的状态，所以记录最新一次的操作，自动扩容使用
	}

	//readOp常量描述了对缓冲区执行的最后一个操作
	type readOp int8

	// Don't use iota for these, as the values need to correspond with the
	// names and comments, which is easier to see when being explicit.
	const (
	   opRead      readOp = -1 // Any other read operation.  任何其他操作
	   opInvalid   readOp = 0  // Non-read operation.    没有读操作
	   opReadRune1 readOp = 1  // Read rune of size 1.  读取大小为1的字符 （由于UTF-8字符可能包含1-4个字节）
	   opReadRune2 readOp = 2  // Read rune of size 2.  读取大小为2的字符
	   opReadRune3 readOp = 3  // Read rune of size 3.  读取大小为3的字符
	   opReadRune4 readOp = 4  // Read rune of size 4.  读取大小为4的字符


	func (b *Buffer) Len() int { return len(b.buf) - b.off }
	作用：返回缓冲区未读部分的字节数


	func (b *Buffer) Truncate(n int) {
	   if n == 0 {
		  b.Reset()
		  return
	   }
	   b.lastRead = opInvalid
	   if n < 0 || n > b.Len() {
		  panic("bytes.Buffer: truncation out of range")
	   }
	   b.buf = b.buf[:b.off+n]
	}
	作用：从缓冲区中丢弃除前n个未读字节以外的所有字节，但继续使用相同的已分配存储。(已读的数据不会删除)
	如果n为负或大于缓冲区的长度，则会发生panic



	func (b *Buffer) Reset() {
	   b.buf = b.buf[:0]
	   b.off = 0
	   b.lastRead = opInvalid
	}
	将缓冲区重置为空，但它会保留底层存储以供将来的写入使用。（清空数据，cap不变）
	offset 偏移量置为0
	lastRead置为未读取

*/
func byteRead() {
	// 1.申明缓冲区大小
	byteSlice := make([]byte, 20)
	byteSlice[0] = 1                                                  // 将缓冲区第一个字节置1
	fmt.Printf("切片len:%d,切片cap:%d\n", len(byteSlice), cap(byteSlice)) // 切片len:20,切片cap:20

	//  2.根据buf初始化buffer缓冲区
	byteBuffer := bytes.NewBuffer(byteSlice)   // 创建20字节缓冲区 len = 20 off = 0
	fmt.Printf("未读len:%d\n", byteBuffer.Len()) //len:20

	// 3.读取元素
	c, _ := byteBuffer.ReadByte()                       // off+=1
	fmt.Printf("未读len:%d, c=%d\n", byteBuffer.Len(), c) // len = 20 off =1   打印len:19, c=1

	// 重置
	byteBuffer.Reset()                         // len = 0 off = 0
	fmt.Printf("未读len:%d\n", byteBuffer.Len()) // 打印len=0

	// 写入元素
	byteBuffer.Write([]byte("hello byte buffer")) // 写缓冲区  len+=17
	fmt.Printf("未读len:%d\n", byteBuffer.Len())    // 打印len=17

	// 跳过元素
	byteBuffer.Next(4)           // 跳过4个字节 off+=4
	c, _ = byteBuffer.ReadByte() // 读第5个字节 off+=1
	fmt.Printf("第5个字节:%d\n", c)  // 打印:111(对应字母o)    len=17 off=5

	// Truncate留下剩下未读的3个字节,使用相同的分配空间,相当于offset+3 = len
	byteBuffer.Truncate(3)                     // 将未字节数置为3        len=off+3=8   off=5
	fmt.Printf("未读len:%d\n", byteBuffer.Len()) // 打印len=3为未读字节数  上面len=8是底层切片长度
	byteBuffer.WriteByte(96)                   // len+=1=9 将y改成A

	byteBuffer.Next(3)           // len=9 off+=3=8
	c, _ = byteBuffer.ReadByte() // off+=1=9    c=96
	fmt.Printf("第9个字节:%d\n", c)  // 打印:96
}

func osRead() {
	//step1：打开文件
	fileName := "chapter01_input_output/danny.txt"
	// os库主要是处理操作系统操作的，它作为Go程序和操作系统交互的桥梁。创建文件、打开或者关闭文件、Socket等等这些操作和都是和操作系统挂钩的，
	//	所以都通过os库来执行
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	//step3：关闭文件
	defer file.Close()

	//step2：读取数据
	bs := make([]byte, 4, 4)
	/*
		//第一次读取
		n,err :=file.Read(bs)
		fmt.Println(err) //<nil>
		fmt.Println(n) //4
		fmt.Println(bs) //[97 98 99 100]  // ASCII值
		fmt.Println(string_test(bs)) //abcd

		//第二次读取
		n,err = file.Read(bs)
		fmt.Println(err)//<nil>
		fmt.Println(n)//4
		fmt.Println(bs) //[101 102 13 10]
		fmt.Println(string_test(bs)) //ef

		//第三次读取
		n,err = file.Read(bs)
		fmt.Println(err) //<nil>
		fmt.Println(n) //4
		fmt.Println(bs) //[65 66 13 10]
		fmt.Println(string_test(bs)) //AB

		//第四次读取
		n,err = file.Read(bs)
		fmt.Println(err) //EOF
		fmt.Println(n) //0
	*/
	n := -1
	for {
		n, err = file.Read(bs)
		// io库属于底层接口定义库。它的作用主要是定义个I/O的基本接口和个基本常量，并解释这些接口的功能
		if n == 0 || err == io.EOF {
			fmt.Println("读取到了文件的末尾，结束读取操作。。")
			break
		}
		fmt.Println(string(bs[:n]))
	}
}
