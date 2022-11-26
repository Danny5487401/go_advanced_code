# 符号 & 符号表

符号表记录了程序中全局函数、全局变量、局部非静态变量与链接器符号解析、重定位相关的信息。

- [代码](chapter31_tool/01_build/02_ldflags/build.go)
```go
package main

import "fmt"

var buildVer string
Ï
func main() {
	fmt.Println("link 传参数为: ", buildVer)

}
```

其中的包名main、函数名main.main、导入的外部包名fmt、引用的外部函数fmt.Println，这些都属于符号的范畴


## 快速查看符号&符号表


go tool compile build.go会输出一个文件build.o，这里的main.o是一个可重定位目标文件，但是其文件格式却不能被readelf、nm分析，
因为它是go自己设计的一种对象文件格式，在proposal: build a better linker种有提及，要分析main.o只能通过go官方提供的工具


```shell
➜  02_ldflags git:(feature/asm) ✗ go tool compile build.go 
➜  02_ldflags git:(feature/asm) ✗ ls
build      build.go   build.o    ldflags.md
➜  02_ldflags git:(feature/asm) ✗ go tool nm build.o      
    15a2 D ""..inittask
    15c2 r ""..stmp_0<1>
    15c2 B "".buildVer
    1592 T "".init
    14e2 T "".main
    15d2 r "".main.stkobj<1>
         U fmt..inittask
         U fmt.Fprintln
    16d1 R gclocals·33cdeccccebe80329f1fdbee7f5874cb
    16bf R gclocals·69c1753bd5f81501d95132d08af04464
    16c7 R gclocals·ef901d0ae51b5399f7d4b5dfa3bc0b42
    16d9 ? go.cuinfo.packagename.
         U go.info.[]interface {}
         U go.info.error
    16dd ? go.info.fmt.Println$abstract
         U go.info.int
    1707 R go.itab.*os.File,io.Writer
    1687 R go.string."link 传参数为: "
         U gofile..$GOROOT/src/fmt/print.go
         U gofile../Users/python/Desktop/go_advanced_code/chapter31_tool/01_build/02_ldflags/build.go
         U os.(*File).Write
         U os.Stdout
    1636 R runtime.gcbits.01
    1637 R runtime.gcbits.02
    1638 R runtime.gcbits.0a
    172f R runtime.memequal64·f
    1727 R runtime.nilinterequal·f
    17bf R type.*[]interface {}
    1737 R type.*interface {}
         U type.*os.File
    16ba R type..importpath.fmt.
    16a9 R type..namedata.*[]interface {}-
    169a R type..namedata.*interface {}-
    17f7 R type.[]interface {}
         U type.int
    176f R type.interface {}
         U type.io.Writer

```

go tool nm和Linux下binutils提供的nm，虽然支持的对象文件格式不同，但是其输出格式还是相同的

- 第一列，symbol value，表示定义符号处的虚拟地址（如变量名对应的变量地址）；
- 第二列，symbol type，用小写字母表示局部符号，大写则为全局符号（uvw例外）；
```shell
"A" The symbol's value is absolute, and will not be changed by further linking.

"B"
"b" （全称为 bss）是未初始化的数据

"C" The symbol is common.  Common symbols are uninitialized data.  When linking, multiple common symbols may appear with the same name.  If the symbol is defined anywhere, the common symbols are treated as undefined references.

"D"
"d"  data section 已初始化的数据

"G"
"g" The symbol is in an initialized data section for small objects.  Some object file formats permit more efficient access to small data objects, such as a global int variable as opposed to a large global array.

"i" For PE format files this indicates that the symbol is in a section specific to the implementation of DLLs.  For ELF format files this indicates that the symbol is an indirect function.  This is a GNU extension to the standard set of ELF symbol types.  It indicates a symbol which if referenced by a relocation does not evaluate to its address, but instead must be invoked at runtime.  The runtime execution will then return the value to be used in the relocation.

"I" The symbol is an indirect reference to another symbol.

"N" The symbol is a debugging symbol.

"p" The symbols is in a stack unwind section.

"R"
"r" The symbol is in a read only data section.

"S"
"s" The symbol is in an uninitialized data section for small objects.

"T"
"t" The symbol is in the text (code) section.

"U" The symbol is undefined.

"u" The symbol is a unique global symbol.  This is a GNU extension to the standard set of ELF symbol bindings.  For such a symbol the dynamic linker will make sure that in the entire process there is just one symbol with this name and type in use.

"V"
"v" The symbol is a weak object.  When a weak defined symbol is linked with a normal defined symbol, the normal defined symbol is used with no error.  When a weak undefined symbol is linked and the symbol is not defined, the value of the weak symbol becomes zero with no error.  On some systems, uppercase indicates that a default value has been specified.

"W"
"w" The symbol is a weak symbol that has not been specifically tagged as a weak object symbol.  When a weak defined symbol is linked with a normal defined symbol, the normal defined symbol is used with no error.  When a weak undefined symbol is linked and the symbol is not defined, the value of the symbol is determined in a system-specific manner without error.  On some systems, uppercase indicates that a default value has been specified.

"-" The symbol is a stabs symbol in an a.out object file.  In this case, the next values printed are the stabs other field, the stabs desc field, and the stab type.  Stabs symbols are used to hold debugging information.

"?" The symbol type is unknown, or object file format specific.
```

- 第三列，symbol name，符号名，对应字符串是存储在字符串表中，由符号表引用；


## 源码

```go
// /Users/python/go/go1.18/src/debug/elf/elf.go

// ELF32 Symbol.
type Sym32 struct {
	Name  uint32
	Value uint32
	Size  uint32
	Info  uint8
	Other uint8
	Shndx uint16
}

```