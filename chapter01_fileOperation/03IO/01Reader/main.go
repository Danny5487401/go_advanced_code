// io 包这些接口和原始的操作以不同的实现包装了低级操作，客户不应假定它们对于并行执行是安全的
// 在io包中最重要的是两个接口：Reader和Writer接口

package main

import (
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
	//读取本地aa.txt文件中的数据
	//step1：打开文件
	fileName := "E:\\go_advanced_code\\chapter01_fileOperation\\danny.txt"
	file,err := os.Open(fileName)
	if err != nil{
		fmt.Println("err:",err)
		return
	}
	//step3：关闭文件
	defer file.Close()

	//step2：读取数据
	bs := make([]byte,4,4)
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
	for{
		n,err = file.Read(bs)
		if n == 0 || err == io.EOF{
			fmt.Println("读取到了文件的末尾，结束读取操作。。")
			break
		}
		fmt.Println(string(bs[:n]))
	}
}