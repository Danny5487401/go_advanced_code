/*
原因：
	io操作本身的效率并不低，低的是频繁的访问本地磁盘的文件。
解决：
	所以bufio就提供了缓冲区(分配一块内存)，读和写都先在缓冲区中，最后再读写文件，来降低访问本地磁盘的次数，从而提高效率。

原理：
	把文件读取进缓冲（内存）之后再读取的时候就可以避免文件系统的io 从而提高速度。
	同理，在进行写操作时，先把文件写入缓冲（内存），然后由缓冲写入文件系统。
	看完以上解释有人可能会表示困惑了，直接把 内容->文件 和 内容->缓冲->文件相比， 缓冲区好像没有起到作用嘛。
	其实缓冲区的设计是为了存储多次的写入，最后一口气把缓冲区内容写入文件
分类：
	主要分三部分Reader、Writer、Scanner,分别是读数据、写数据和扫描器三种数据类型

bufio 封装了io.Reader或io.Writer接口对象，并创建另一个也实现了该接口的对象

io.Reader或io.Writer 接口实现read() 和 write() 方法，对于实现这个接口的对象都是可以使用这两个方法的
*/

package main

import (
	"bufio"
	"fmt"
	_ "io"
	"os"
	"strings"
	"time"
)

func main() {
	/*
		bufio:高效io读写
			buffer缓存
			io：input/output

		将io包下的Reader，Write对象进行包装，带缓存的包装，提高读写的效率

			ReadBytes()
			ReadString()
			ReadLine()

	*/

	fileName := "./chapter01_fileOperation/dannyBufio.txt"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//创建Reader对象
	//b1 := bufio.NewReader(file) //  底层调用NewReaderSize(rd, defaultBufSize)
	b1 := bufio.NewReaderSize(file, 200) //  底层调用NewReaderSize(rd, defaultBufSize)

	//// 1. Read()高效读取
	//p := make([]byte, 512)
	//n1, err := b1.Read(p) // 单次读取
	//fmt.Println(n1)
	//fmt.Println(string(p[:n1]))

	// 读取方式一：readSlice
	//for {
	//	// ReadSlice返回的 []byte 是指向 Reader 中的 buffer，而不是 copy 一份返回，也正因为如此，通常我们会使用 ReadBytes 或 ReadString
	//	s1, err := b1.ReadSlice('\n') //注意是单引号
	//	/*
	//		如果 ReadSlice 在找到界定符之前遇到了 error ，它就会返回缓存中所有的数据和错误本身（经常是 io.EOF）。
	//		如果在找到界定符之前缓存已经满了，ReadSlice 会返回 bufio.ErrBufferFull 错误。
	//		当且仅当返回的结果（line）没有以界定符结束的时候，ReadSlice 返回err != nil，也就是说，如果ReadSlice 返回的结果 line 不是以界定符 delim 结尾,
	//		那么返回的 err也一定不等于 nil（可能是bufio.ErrBufferFull或io.EOF
	//
	//	*/
	//	if err == io.EOF {
	//		fmt.Println("读取完毕。。")
	//		break
	//	}
	//	// 注意:这里的界定符可以是任意的字符，可以将上面代码中的'\n'改为'm'试试。同时，返回的结果是包含界定符本身的，
	//	// 输出结果有一空行就是’\n’本身(line携带一个’\n’,printf又追加了一个’\n’)。
	//	fmt.Printf(string(s1))
	//}

	// 后三个方法最终都是调用ReadSlice来实现的

	////读取方式二：ReadLine() 读一行
	/*
		ReadLine 尝试返回单独的行，不包括行尾的换行符。如果一行大于缓存，isPrefix 会被设置为 true，同时返回该行的开始部分（等于缓存大小的部分）。
		该行剩余的部分就会在下次调用的时候返回。当下次调用返回该行剩余部分时，isPrefix 将会是 false 。跟 ReadSlice 一样，返回的 line 只是 buffer 的引用，
		在下次执行IO操作时，line 会无效。可以将 ReadSlice 中的例子该为 ReadLine 试试。

	*/
	//data, flag, err := b1.ReadLine()
	//fmt.Println(flag)
	//fmt.Println(err)
	//fmt.Println(data) //byte数据
	//fmt.Println(string(data))

	//读取方式三：ReadBytes() 读到分隔符
	//data,err :=b1.ReadBytes('\n')
	//fmt.Println(err)
	//fmt.Println(string_test(data))

	//读取方式四：ReadString(),底层调用readBytes
	s1, err := b1.ReadString('\n')
	fmt.Println(err)
	fmt.Println(s1)

	// 全部读取
	//for {
	//	s1, err := b1.ReadString('\n')
	//	if err == io.EOF {
	//		fmt.Println("读取完毕。。")
	//		break
	//	}
	//	fmt.Printf("%#v\n", s1)
	//}

	// peek 该方法只是“窥探”一下 Reader 中没有读取的 n 个字节。好比栈数据结构中的取栈顶元素，但不出栈。
	reader := bufio.NewReaderSize(strings.NewReader("http://studygolang.com.\t It is the home of gophers"), 14)
	go Peek(reader)
	go reader.ReadBytes('\t')
	time.Sleep(1e8)

	//Scanner
	input := "foo   bar      baz"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords) //  bufio.ScanWords返回通过“空格”分词的单词
	//scanner.Split(bufio.ScanLines)
	//scanner.Split(bufio.ScanRunes)
	//scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

// 不安全
/*
ReadSlice一样，返回的 []byte 只是 buffer 中的引用，在下次IO操作后会无效，可见该方法（以及ReadSlice这样的，返回buffer引用的方法）对多 goroutine 是不安全的，
也就是在多并发环境下，不能依赖其结果
*/
func Peek(reader *bufio.Reader) {
	line, _ := reader.Peek(14)
	fmt.Printf("%s\n", line)
	//time.Sleep(1) // 注释
	fmt.Printf("%s\n", line)

	/*结果
	如果我们将例子中注释掉的 time.Sleep(1) 取消注释（这样调度其他 goroutine 执行），再次运行
	*/
}
