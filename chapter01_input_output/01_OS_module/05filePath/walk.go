package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	pwd, _ := os.Getwd()

	// 遍历指定目录(包括子目录)，对遍历的项目用walkFn函数进行处理
	filepath.Walk(pwd, WalkFunc)

}

/*
	规则如下：

	文件处理函数定义如下，如果 WalkFunc 返回 nil，则 Walk 函数继续遍历，如果返回 SkipDir，则 Walk 函数会跳过当前目录（如果当前遍历到的是文件，
		则同时跳过后续文件及子目录），继续遍历下一个目录。如果返回其它错误，则 Walk 函数会中止遍历。在 Walk 遍历过程中，如果遇到错误，
		则会将错误通过 err 传递给WalkFunc 函数，同时 Walk 会跳过出错的项目，继续处理后续项目

*/
func WalkFunc(fpath string, info os.FileInfo, err error) error {
	// 获取path中最后一个分隔符之后的部分(不包含分隔符)
	base := filepath.Base(fpath)

	// 匹配当前目录下名字为3个字符的
	var match bool
	if match, err = filepath.Match("???", base); match {
		fmt.Println("File:", fpath, "IsDir:", info.IsDir(), "size:", info.Size())
		return nil
	}
	return err
}

/*
match函数根据pattern来判断name是否匹配，如果匹配则返回true，pattern 规则如下：

       可以使用 ? 匹配单个任意字符（不匹配路径分隔符）。

       可以使用 * 匹配 0 个或多个任意字符（不匹配路径分隔符）。

       可以使用 [] 匹配范围内的任意一个字符（可以包含路径分隔符）。

       可以使用 [^] 匹配范围外的任意一个字符（无需包含路径分隔符）。

       [] 之内可以使用 - 表示一个区间，比如 [a-z] 表示 a-z 之间的任意一个字符。

       反斜线用来匹配实际的字符，比如 \* 匹配 *，\[ 匹配 [，\a 匹配 a 等等。

       [] 之内可以直接使用 [ * ?，但不能直接使用 ] -，需要用 \]、\- 进行转义。

*/
