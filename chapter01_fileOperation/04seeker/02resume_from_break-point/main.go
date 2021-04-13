//断点续传
//Q1：如果你要传的文件，比较大，那么是否有方法可以缩短耗时？
//Q2：如果在文件传递过程中，程序因各种原因被迫中断了，那么下次再重启时，文件是否还需要重头开始？
//Q3：传递文件的时候，支持暂停和恢复么？即使这两个操作分布在程序进程被杀前后

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	/*
		断点续传：
			文件传递：文件复制
				E:\\go_advanced_code\\chapter01_fileOperation\\resume_from_breakingPoint.jpg"

			复制到
				E:\go_advanced_code\chapter01_fileOperation\dest.jpg

		思路：
			边复制，边记录复制的总量
	*/

	srcFile := "E:\\go_advanced_code\\chapter01_fileOperation\\resume_from_breakingPoint.jpg"
	destFile := "E:\\go_advanced_code\\chapter01_fileOperation\\dest.jpg"
	tempFile := destFile + "temp.txt"
	//fmt.Println(tempFile)
	file1, _ := os.Open(srcFile)
	file2, _ := os.OpenFile(destFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	file3, _ := os.OpenFile(tempFile, os.O_CREATE|os.O_RDWR, os.ModePerm)

	defer file1.Close()
	defer file2.Close()
	//1.读取临时文件中的数据，根据seek
	file3.Seek(0, io.SeekStart)
	bs := make([]byte, 100, 100)
	n1, err := file3.Read(bs)
	fmt.Println(n1)
	countStr := string(bs[:n1])
	fmt.Println(countStr)
	//count,_:=strconv.Atoi(countStr)
	count, _ := strconv.ParseInt(countStr, 10, 64)
	fmt.Println(count)

	//2. 设置读，写的偏移量
	file1.Seek(count, 0)
	file2.Seek(count, 0)
	data := make([]byte, 1024, 1024)
	n2 := -1            // 读取的数据量
	n3 := -1            //写出的数据量
	total := int(count) //读取的总量

	for {
		//3.读取数据
		n2, err = file1.Read(data)
		if err == io.EOF {
			fmt.Println("文件复制完毕。。")
			file3.Close()
			os.Remove(tempFile)
			break
		}
		//将数据写入到目标文件
		n3, _ = file2.Write(data[:n2])
		total += n3
		//将复制总量，存储到临时文件中
		file3.Seek(0, io.SeekStart)
		file3.WriteString(strconv.Itoa(total))

		//假装断电
		//if total>8000{
		//	panic("假装断电了。。。，假装的。。。")
		//}
	}

}