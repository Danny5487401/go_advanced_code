/*  编译
注：先执行
gcc -c hello.c   // -c 只编译不链接

// ar 工具：为了创建和管理连接程序使用的目标库文档
ar -cru libhello.a hello.o    // 将一组编译过的文件合并为一个文件  -c 无提示模式创建文件包, -r 在文件包中代替文件,-u 与r共同使用,用来仅取代那些在生成文件包之后改动过的文件.

 */
package main

//使用#cgo定义库路径
// 下面是相对路径.
// 连接选项:可以使用#cgo扩展预处理命令指定CFLAGS/LDFLAGS。

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lhello
#include "hello.h"
*/
import "C"

func main() {
	C.hello()
}
