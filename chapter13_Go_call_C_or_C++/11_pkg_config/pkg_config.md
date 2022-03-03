# Golang使用pkg-config自动获取头文件和链接库的方法

在Golang中使用cgo调用C库的时候，如果需要引用很多不同的第三方库，那么使用#cgo CFLAGS:和#cgo LDFLAGS:的方式会引入很多行代码。

首先这会导致代码很丑陋，最重要的是如果引用的不是标准库，头文件路径和库文件路径写死的话就会很麻烦。
一旦第 三方库的安装路径变化了，Golang的代码也要跟着变化，所以使用pkg-config无疑是一种更为优雅的方法，

## 使用
在路径 /Users/xiaxin/Desktop/go_advanced_code/chapter13_Go_call_C_or_C++/11_pkg_config
下安装了一个名称为hello的第三方C语言库，其目录结构如下所示，在hello_world.h中只定义了一个接口函数hello，该函数接收一个char *字符串作为变量并调用printf将其打印到标准输出