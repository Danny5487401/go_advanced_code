package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main()  {
	pwd,_ := os.Getwd()
	/*
	规则如下：

	文件处理函数定义如下，如果 WalkFunc 返回 nil，则 Walk 函数继续遍历，如果返回 SkipDir，则 Walk 函数会跳过当前目录（如果当前遍历到的是文件，
		则同时跳过后续文件及子目录），继续遍历下一个目录。如果返回其它错误，则 Walk 函数会中止遍历。在 Walk 遍历过程中，如果遇到错误，
		则会将错误通过 err 传递给WalkFunc 函数，同时 Walk 会跳过出错的项目，继续处理后续项目

	 */

	filepath.Walk(pwd,func(fpath string, info os.FileInfo, err error) error {
		if match,err := filepath.Match("???",filepath.Base(fpath)); match {
			fmt.Println("path:",fpath)
			//fmt.Println("info:",info)
			return err
			 }
		return nil
		})


	
}
