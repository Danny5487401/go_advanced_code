/*
文件操作：file类是在os包中的，封装了底层的文件描述符和相关信息，同时封装了Read和Write的实现。
*/

package main

import (
	"fmt"
	"os"
)

func main() {
	/*
		FileInfo：文件信息
			interface
				Name()，文件名
				Size()，文件大小，字节为单位
				IsDir()，是否是目录
				ModTime()，修改时间
				Mode()，权限

	*/
	// 获取文件的信息，里面有文件的名称，大小，修改时间等
	fileInfo, err := os.Stat("chapter01_fileOperation/danny.txt")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	fmt.Printf("%T\n", fileInfo)
	//文件名
	fmt.Println(fileInfo.Name())
	//文件大小
	fmt.Println(fileInfo.Size())
	//是否是目录
	fmt.Println(fileInfo.IsDir()) //IsDirectory
	//修改时间
	fmt.Println(fileInfo.ModTime())
	//权限
	fmt.Println(fileInfo.Mode()) //-rw-r--r--
}
