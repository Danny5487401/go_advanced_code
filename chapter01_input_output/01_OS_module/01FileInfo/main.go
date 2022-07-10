/*
文件操作：file类是在os包中的，封装了底层的文件描述符和相关信息，同时封装了Read和Write的实现。
*/

package main

import (
	"fmt"
	"os"
)

/*
type FileInfo interface {
    Name() string       // 文件的名字（不含扩展名）
    Size() int64        // 普通文件返回值表示其大小,单位byte字节；其他文件的返回值含义各系统不同
    Mode() FileMode     // 文件的模式位
    ModTime() time.Time // 文件的修改时间
    IsDir() bool        // 等价于Mode().IsDir()
    Sys() interface{}   // 底层数据来源（可以返回nil）
}

*/

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
	fmt.Printf("%T\n", fileInfo) // *os.fileStat 具体实现

	//文件名
	fmt.Println(fileInfo.Name()) // danny.txt
	//文件大小
	fmt.Println(fileInfo.Size()) // 10
	//是否是目录
	fmt.Println(fileInfo.IsDir()) //IsDirectory
	//修改时间
	fmt.Println(fileInfo.ModTime())
	//权限
	fmt.Println(fileInfo.Mode()) //-rw-r--r--
}
