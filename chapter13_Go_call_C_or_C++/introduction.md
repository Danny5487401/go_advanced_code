# go 命令使用 cgo  
为了使用 cgo，你需要在普通的 Go 代码中导入一个伪包 "C"。这样 Go 代码就可以引用一些 C 的类型 (如 C.size_t)、变量 (如 C.stdout)、或函数 (如 C.putchar)
如果对 "C" 的导入语句之前紧贴着是一段注释，那么这段注释被称为前言，它被用作编译 C 部分的头文件。如下面例子所示：
```shell script
// #include <stdio.h>
// #include <errno.h>
import "C"
```
前言中可以包含任意 C 代码，包括函数和变量的声明和定义。虽然他们是在 "C" 包里定义的，但是在 Go 代码里面依然可以访问它们。所有在前言中声明的名字都可以被 Go 代码使用，即使名字的首字母是小写的。static 变量是个例外：它不能在 Go 代码中被访问。但是 static 函数可以在 Go 代码中访问
# 标准库案例
$GOROOT/misc/cgo/stdio 和 $GOROOT/misc/cgo/gmp
多个指令定义的值会被串联到一起。这些指令可以包括一系列构建约束，用以限制对满足其中一个约束的系统的影响
```shell script
// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo amd64 386 CFLAGS: -DX86=1
// #cgo LDFLAGS: -lpng
// #include <png.h>
import "C"
```
CPPFLAGS 和 LDFLAGS 也可以通过 #cgo pkg-config 命令来通过 pkg-config 来获取。随后的指令指定获取的包名
```go
// #cgo: pkg-config: png cairo
// #include <png.h>
import "C"
```
默认的 pkg-config 工具配置可以通过设置 PKG_CONFIG 环境变量来修改
处于安全原因，只有一部分标志 (flag) 允许设置，特别是 -D，-I，以及 -l。如果想允许设置额外的 flag，可以设置 CGO_CFLAGS_ALLOW，按正则条件匹配新的 flag。如果想禁止当前允许的 flag，设置 CGO_CFLAGS_DISALLOW，按正则匹配被禁止的指令。在这两种情况下，正则匹配都必须是完全匹配：如果想允许 -mfoo=bar 指令，设置 CGO_CFLAGS_ALLOW='-mfoo.*'，而不能仅仅设置 CGO_CFLAGS_ALLOW='-mfoo'。类似名称的变量控制 CPPFLAGS, CXXFLAGS, FFLAGS, 以及 LDFLAGS。

当构建的时候，CGO_CFLAGS, CGO_CPPFLAGS, CGO_CXXFLAGS, CGO_FFLAGS 和 CGO_LDFLAGS 环境变量会被赋值为由上述指令派生的值。跟包有关的 flag 应该使用指令来设置，而不是通过环境变量来设置，这样构建工作可以处于未修改的环境中。从环境变量中获取的值不属于上述描述的安全限制。
编译事项： 
在一个 Go 包中的所有 CPPFLAGS 和 CFLAGS 指令都会被串联到一起，然后被用来编译这个包中的 C 文件。包中所有的 CPPFLAGS 和 CXXFLAGS 指令都会被串联到一起，用于编译包中的 C++ 文件。包中所有的 CPPFLAGS 和 FFLAGS 指令都会被串联到一起，用来编译包中的 Fortran 文件。一个程序中所有包中的 LDFLAGS 指令都会被串联到一起，并在连接 (link) 的时候使用。所有的 pkg-config 指令都会被串联到一起，并同时发给 pkg-config，来添加每个合适的 flag 命令行的集合
在解析 cgo 指令时，所有字符串中的 ${SRCDIR} 都会被替换成当前源文件所在的绝对路径。这样允许提前编译的静态库被包含在包路径中且被正确地链接。比如，如果包 foo 在 /go/src/foo 路径下：
```go
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
//会被扩展成：
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo

```

当 Go tools 发现一个或多个 Go 文件使用特殊的引用 "C" 时，它会寻找当前路径中的非 Go 文件，并把这些文件编译为 Go 包的一部分。任意 .c, .s, 或 .S 文件都会被 C 编译器编译。任意 .cc, .cpp, 或 .cxx 文件都会被 C++ 编译器编译。任意 .f, .F, .for 或 .f90 文件都会被 fortran 编译器编译。任意 .h, .hh, .hpp 或 .hxx 文件都不会被分别编译，但是如果这些头文件被修改了，那么 C 和 C++ 文件会被重新编译。默认的 C 和 C++ 编译器都可以分别通过设置 CC 和 CXX 环境变量来修改。这些环境变量可能包括命令行选项。

当在被期望的系统上构建 Go 时，cgo tool 默认是开启的。在交叉编译时，它默认是关闭的。你可以通过设置 CGO_ENABLED 环境变量来控制开启和关闭：设置为 1 表示启用 cgo，设置为 0 表示不启用 cgo。如果 cgo 被启用，则 go tool 会设置构建约束“cgo”。

在交叉编译时，你必须为 cgo 指定一个 C 交叉编译器。你可以在使用 make 构建工具链 (toolchain) 时设置通用的 CC_FOR_TARGET 或更明确的 CC_FOR_${GOOS}_${GOARCH} (比如 CC_FOR_linux_arm) 环境变量，或者在任意运行 go 工具时设置 CC 环境变量。

CXX_FOR_TARGET, CXX_FOR_${GOOS}_${GOARCH} 以及 CXX 环境变量使用方式类似
<table border="1" cellpadding="1" cellspacing="1" style="width:500px;"><tbody><tr><th> <p>C 类型名称</p> </th><th>Go 类型名称</th></tr></tbody><tbody><tr><td>char</td><td>C.char</td></tr><tr><td>signed char</td><td>C.schar</td></tr><tr><td>unsigned char</td><td>C.uchar</td></tr><tr><td>short</td><td>C.short</td></tr><tr><td>unsighed short</td><td>C.ushort</td></tr><tr><td>int</td><td>C.int</td></tr><tr><td>unsigned int</td><td>C.uint</td></tr><tr><td>long</td><td>C.long</td></tr><tr><td>unsigned long</td><td>C.ulong</td></tr><tr><td>long long</td><td>C.longlong</td></tr><tr><td>unsigned long long</td><td>C.ulonglong</td></tr><tr><td>float</td><td>C.float</td></tr><tr><td>double</td><td>C.double</td></tr><tr><td>complex float</td><td>C.complexfloat</td></tr><tr><td>complex double</td><td>C.complexdouble</td></tr><tr><td>void*</td><td>unsafe.Pointer</td></tr><tr><td>__int128_t&nbsp; &nbsp;__uint128_t</td><td>[16]byte</td></tr></tbody></table>
一些通常在 Go 中被表示为指针类型的特殊 C 类型会被表示成 uintptr。下面的特殊场景会对此进行介绍。

当直接访问 C 中的结构体、联合、或枚举类型时，在名字前面加上 struct_、union_、或 enum_，就像 C.struct_stat 这样。

C 的任意类型 T 的大小，在 Go 中用 C.sizeof_T 表示，比如 C.sizeof_struct_stat。

可以在 Go 文件中声明一个带有特殊类型 _GoString_ 类型的 C 函数。可以使用普通的 Go 字符串调用这个函数。可以通过调用这些 C 函数来获取字符串长度，或指向字符串的指针
```go

size_t _GoStringLen(_GoString_ s);
const char *_GoStringPtr(_GoString_ s);
```
这些函数只能被写在 Go 文件的前言里，而不能写在其他的 C 文件中。C 代码必须不能修改 _GoStringPtr 返回的指针内容。注意字符串内容可能不是 NULL结尾。

由于通常情况下 Go 不支持 C 的联合类型，C 的联合类型在 Go 中被表示为相同长度的比特数组。

Go 的结构体不能嵌入 C 的类型。

对于一个非空的 C 结构体，如果它结尾的字段大小为 0，那么 Go 代码无法引用这个字段。为了获取到这样字段的地址，你只能先获取结构体的地址，然后将地址加上这个结构体的大小。这也是能获取到这个字段的唯一方式。

Cgo 把 C 类型转换成等价的不可输出的 Go 类型。因此 Go 包不应该在它的输出接口中暴露 C 类型：同一个 C 类型，在不同包里是不一样的。

任意 C 函数 (即使是 void 函数) 可能会被在多赋值场景下被调用，以获取返回值 (如果有的话)，和 C errno 值 (作为 error 值)。如果 C 函数返回类型是 void，那么相应的值可以用 _ 代替。
```go

n, err = C.sqrt(-1)
_, err := C.voidFunc()
var n, err = C.sqrt(1)
```