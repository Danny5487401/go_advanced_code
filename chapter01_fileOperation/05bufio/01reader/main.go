//io操作本身的效率并不低，低的是频繁的访问本地磁盘的文件。
//所以bufio就提供了缓冲区(分配一块内存)，读和写都先在缓冲区中，最后再读写文件，来降低访问本地磁盘的次数，从而提高效率。

/*
把文件读取进缓冲（内存）之后再读取的时候就可以避免文件系统的io 从而提高速度。
同理，在进行写操作时，先把文件写入缓冲（内存），然后由缓冲写入文件系统。
看完以上解释有人可能会表示困惑了，直接把 内容->文件 和 内容->缓冲->文件相比， 缓冲区好像没有起到作用嘛。
其实缓冲区的设计是为了存储多次的写入，最后一口气把缓冲区内容写入文件
 */

package main

import (
	"fmt"
	"os"
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

	fileName:="E:\\go_advanced_code\\chapter01_fileOperation\\dannyBufio.txt"
	file,err := os.Open(fileName)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer file.Close()

	//创建Reader对象
	//b1 := bufio.NewReader(file)

	//// 1. Read()高效读取
	//p := make([]byte,1024)
	//n1,err := b1.Read(p)
	//fmt.Println(n1)
	//fmt.Println(string_test(p[:n1]))

	////2.ReadLine()
	//data,flag,err := b1.ReadLine()
	//fmt.Println(flag)
	//fmt.Println(err)
	//fmt.Println(data)
	//fmt.Println(string_test(data))

	//3.ReadString() 读到分隔符
	//s1,err :=b1.ReadString('\n')
	//fmt.Println(err)
	//fmt.Println(s1)
	//
	// s1,err = b1.ReadString('\n')
	// fmt.Println(err)
	// fmt.Println(s1)
	//
	//s1,err = b1.ReadString('\n')
	//fmt.Println(err)
	//fmt.Println(s1)
	//
	//for{
	//	s1,err := b1.ReadString('\n')
	//	if err == io.EOF{
	//		fmt.Println("读取完毕。。")
	//		break
	//	}
	//	fmt.Println(s1)
	//}

	//4.ReadBytes() 读到分隔符
	//data,err :=b1.ReadBytes('\n')
	//fmt.Println(err)
	//fmt.Println(string_test(data))


	//Scanner
	//s2 := ""
	//fmt.Scanln(&s2)
	//fmt.Println(s2)

	//b2 := bufio.NewReader(os.Stdin)
	//s2, _ := b2.ReadString('\n')
	//fmt.Println(s2)

}