/*
文件操作：file类是在os包中的，封装了底层的文件描述符和相关信息，同时封装了Read和Write的实现。
*/

package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {

	// 获取文件的信息，里面有文件的名称，大小，修改时间等
	fileInfo, err := os.Stat("chapter01_input_output/files/danny.txt")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	/*
		// A fileStat is the implementation of FileInfo returned by Stat and Lstat.
		type fileStat struct {
			name    string
			size    int64
			mode    FileMode
			modTime time.Time
			sys     syscall.Stat_t
		}
	*/
	fmt.Printf("类型%T\n", fileInfo) // *os.fileStat 具体实现

	//文件名
	fmt.Println(fileInfo.Name()) // danny.txt
	//文件大小
	fmt.Println(fileInfo.Size()) // 10
	//是否是目录
	fmt.Println(fileInfo.IsDir()) //false
	//修改时间
	fmt.Println("对文件上次修改时间", fileInfo.ModTime()) // 等价：fs.modTime = timespecToTime(fs.sys.Mtimespec)
	//权限
	fmt.Println(fileInfo.Mode()) //-rw-r--r--
	// sys
	sys := fileInfo.Sys()
	stat := sys.(*syscall.Stat_t)
	fmt.Println("对文件上次访问时间", time.Unix(stat.Atim.Unix()))
}
