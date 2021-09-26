package main

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
	文件操作：
	1.路径：
		func Getwd() (dir string, err error) // 获取当前工作目录的根路径
		func Chdir(dir string) error // 将工作目录修改为dir
		相对路径：relative
			ab.txt
			相对于当前工程
		绝对路径：absolute
			E:\go_advanced_code\chapter01_fileOperation\danny.txt

		.当前目录
		..上一层
	2.创建文件夹，如果文件夹存在，创建失败
		func Mkdir(name string, perm FileMode) error // 使用指定的权限和名称创建一个文件夹（对于linux的mkdir命令
		func MkdirAll(path string, perm FileMode) error // 使用指定的权限和名称创建一个文件夹，并自动创建父级目录（对于linux的mkdir -p目录

	3.创建文件，Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）
		func Create(name string) (file *File, err error) // 创建一个空文件，注意当文件已经存在时，会直接覆盖掉原文件，不会报错


	4.打开文件：让当前的程序，和指定的文件之间建立一个连接
		func Open(name string) (file *File, err error) // 打开一个文件,注意打开的文件只能读，不能写
		func OpenFile(name string, flag int, perm FileMode) (file *File, err error) // 以指定的权限打开文件

	5.关闭文件：程序和文件之间的链接断开。
		file.Close()

	5.删除文件或目录：慎用，慎用，再慎用
		func Remove(name string) error // 删除指定的文件夹或者目录,不能递归删除，只能删除一个空文件夹或一个文件（对应linux的 rm命令
		func RemoveAll(path string) error // 递归删除文件夹或者文件（对应linux的rm -rf命令）
*/
func main() {

	//1.路径
	wd, _ := os.Getwd()
	fmt.Println("获取当前工作目录的根路径:", wd)

	fileName1 := "chapter01_fileOperation\\danny.txt"
	fileName2 := "./chapter01_fileOperation/danny.txt"
	fmt.Println(filepath.IsAbs(fileName1)) //true
	fmt.Println(filepath.IsAbs(fileName2)) //false
	fmt.Println(filepath.Abs(fileName1))
	fmt.Println(filepath.Abs(fileName2)) // /Users/ruby/go/src/l_file/bb.txt

	//2.创建目录
	//err := os.Mkdir("E:\\go_advanced_code\\chapter01_fileOperation\\testDir",os.ModePerm)
	//if err != nil{
	//	fmt.Println("err:",err)
	//	return
	//}
	//fmt.Println("文件夹创建成功。。")
	//err =os.MkdirAll("E:\\go_advanced_code\\chapter01_fileOperation\\testDir\\aa\\cc\\dd\\ee",os.ModePerm)
	//if err != nil{
	//	fmt.Println("err:",err)
	//	return
	//}
	//fmt.Println("多层文件夹创建成功")

	//3.创建文件:Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）
	//file1,err :=os.Create("E:\\go_advanced_code\\chapter01_fileOperation\\danny2.txt")
	//if err != nil{
	//	fmt.Println("err：",err)
	//	return
	//}
	//fmt.Println(file1)

	//file2,err := os.Create("./chapter01_fileOperation/danny3.txt")//创建相对路径的文件，是以当前工程为参照的
	//if err != nil{
	//	fmt.Println("err :",err)
	//	return
	//}
	//fmt.Println(file2)

	//4.打开文件：
	file3, err := os.Open("E:\\go_advanced_code\\chapter01_fileOperation\\danny2.txt") //只读的
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(file3)
	/* os.OpenFile
	第一个参数：文件名称
	第二个参数：文件的打开方式
		const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
		O_RDONLY int = syscall.O_RDONLY // open the file read-only.
		O_WRONLY int = syscall.O_WRONLY // open the file write-only.
		O_RDWR   int = syscall.O_RDWR   // open the file read-write.
		// The remaining values may be or'ed in to control behavior.
		O_APPEND int = syscall.O_APPEND // append data to the file when writing.
		O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
		O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
		O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
		O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
	)
	第三个参数：文件的权限：文件不存在创建文件，需要指定权限
	*/
	file4, err := os.OpenFile("chapter01_input_output/danny2.txt", os.O_RDONLY|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(file4)
	//
	////5关闭文件，
	file4.Close()

	//6.删除文件或文件夹：
	//删除文件
	err = os.Remove("chapter01_fileOperation\\danny3.txt")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("删除文件成功。。")
	// 删除目录 ：最后一个目录
	//err :=  os.RemoveAll("E:\\go_advanced_code\\chapter01_fileOperation\\testDir\\aa\\cc\\dd\\ee")
	//if err != nil{
	//	fmt.Println("err:",err)
	//	return
	//}
	//fmt.Println("删除目录成功。。")
}
