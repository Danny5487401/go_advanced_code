# 内存对齐

> 现代计算机中内存空间都是按照字节(byte)进行划分的，所以从理论上讲对于任何类型的变量访问都可以从任意地址开始，
> 但是在实际情况中，在访问特定类型变量的时候经常在特定的内存地址访问，所以这就需要把各种类型数据按照一定的规则在空间上排列，而不是按照顺序一个接一个的排放，这种就称为内存对齐，内存对齐是指首地址对齐，而不是说每个变量大小对齐

因为 CPU 访问内存时，并不是逐个字节访问，而是以字（word）为单位访问。
比如 64位CPU的字长（word size）为8bytes，那么CPU访问内存的单位也是8字节，每次加载的内存数据也是固定的若干字长，如8words（64bytes）、16words(128bytes）等



## 编译器为什么要做内存对齐这种事情呢? 

1. 有些CPU可以访问任意地址上的任意数据，而有些CPU只能在特定地址访问数据，因此不同硬件平台具有差异性，这样的代码就不具有移植性，如果在编译时，将分配的内存进行对齐，这就具有平台可以移植性了
2. CPU每次寻址都是要消费时间的，并且CPU 访问内存时，并不是逐个字节访问，而是以字长（word size）为单位访问，所以数据结构应该尽可能地在自然边界上对齐. 
   如果访问未对齐的内存，处理器需要做两次内存访问，而对齐的内存访问仅需要一次访问，内存对齐后可以提升性能

举个例子, 如果不做内存对齐, 那么下面这个结构体的内存分布为
```go
type Test struct {
 	b bool   // 1 byte
 	i3 int32 // 4 bytes
}
```
内存布局：ABBB B

假设当前CPU是32位的，并且没有内存对齐机制，数据可以任意存放，现在有一个int32变量占4byte，存放地址在0x00000002 - 0x00000005(纯假设地址，莫当真)，
这种情况下，每次取4字节的CPU第一次取到[0x00000000 - 0x00000003]，只得到变量1/2的数据，所以还需要取第二次，为了得到一个int32类型的变量，需要访问两次内存并做拼接处理，影响性能。
如果有内存对齐了，int32类型数据就会按照对齐规则在内存中，上面这个例子就会存在地址0x00000000处开始，那么处理器在取数据时一次性就能将数据读出来了，而且不需要做额外的操作，使用空间换时间，提高了效率。

## 对齐系数
每个特定平台上的编译器都有自己的默认"对齐系数"，常用平台默认对齐系数如下：

- 32位系统对齐系数是4
- 64位系统对齐系数是8

这只是默认对齐系数，实际上对齐系数我们是可以修改的，之前写C语言的朋友知道，可以通过预编译指令#pragma pack(n)来修改对齐系数，因为C语言是预处理器的，但是在Go语言中没有预处理器，只能通过tags和命名约定来让Go的包可以管理不同平台的代码，


## go 结构体的内存布局
内存对齐，大家都喜欢拿结构体的内存对齐来举例子，这里要提醒大家一下，不要混淆了一个概念，其他类型也都是要内存对齐的，只不过拿结构体来举例子能更好的理解内存对齐，

GO编译器在编译的时候, 为了保证内存对齐, 对每一个数据类型都给出了对齐保证, 将未对齐的内存留空. 
如果一个类型的对齐保证是4B, 那么其数据存放的起始地址偏移量必是4B 的整数倍. 而编译器给出的这个对齐保证是多少呢? 

规则
> 1 对于结构体的各个成员，第一个成员位于偏移为0的位置，结构体第一个成员的偏移量(offset)为0，以后每个成员相对于结构体首地址的offset都是该成员大小与有效对齐值中较小那个的整数倍，如有需要编译器会在成员之间加上填充字节。

> 2 除了结构成员需要对齐，结构本身也需要对齐，结构的长度必须是编译器默认的对齐长度和成员中最长类型中最小的数据大小的倍数对齐


可以通过内置unsafe包的Sizeof函数来获取一个变量的大小，此外我们还可以通过内置unsafe包的Alignof函数来获取一个变量的对齐系数(不同版本不同平台的编译器不尽相同)

```go
// 结构体变量b1的对齐系数
fmt.Println(unsafe.Alignof(b1))   // 8
// b1每一个字段的对齐系数
fmt.Println(unsafe.Alignof(b1.x)) // 4：表示此字段须按4的倍数对齐
fmt.Println(unsafe.Alignof(b1.y)) // 8：表示此字段须按8的倍数对齐
fmt.Println(unsafe.Alignof(b1.z)) // 1：表示此字段须按1的倍数对齐
```


```go
package main

import (
	"unsafe"
	"fmt"
)

type W struct {
	b byte  // 1 byte
	i int32  // 4 byte
	j int64 // 8byte
}
// b是byte类型，占1个字节；i是int32类型，占4个字节；j是int64类型，占8个字节: 1+4+8=13

func main() {
	printWAlign()
}


func printWAlign() {
    var w = new(W)
    fmt.Printf("结构体 W 的size=%d\n", unsafe.Sizeof(*w))     // size=16
    fmt.Printf("w.b_alignSize=%d\n", unsafe.Alignof(w.b)) // alignSize=1
    fmt.Printf("w.i_alignSize=%d\n", unsafe.Sizeof(w.i))  // alignSize=4
    fmt.Printf("w.j_alignSize=%d\n", unsafe.Sizeof(w.j))  // alignSize=8
}
```
13 !=16
原因：发生了对齐。

对齐值最小是1，这是因为存储单元是以字节为单位。
所以b就在w的首地址，而i的对齐值是4，它的存储地址必须是4的倍数，因此，在b和i的中间有3个填充，同理j也需要对齐，
但因为i和j之间不需要填充，所以w的Sizeof值应该是13+3=16。


在struct中，它的对齐值是它的成员中的最大对齐值。每个成员类型都有它的对齐值，可以用unsafe.Alignof方法来计算。
比如unsafe.Alignof(w.b)就可以得到b在w中的对齐值。
同理，我们可以计算出w.b的对齐值是1，w.i的对齐值是4，w.j的对齐值也是4


## 空结构体字段对齐
Go语言中空结构体的大小为0，如果一个结构体中包含空结构体类型的字段时，通常是不需要进行内存对齐的，




## 检测工具 go vet检查

.golangci.yml 配置
```yaml
linters-settings:
  govet:
    enable-all: true

```


比如下面这个原始结构体
```go
// 处理人信息
type DealerInfo struct {
	DealId2DealerName map[string]string
	DealId2Weight     map[string]int64
	DealerIdList      []string 
	FakePeople        string   
	ManagerId         string   
}

```

golangci-lint 会对结构体的内存字节对齐进行检查，会报下面的错误：
fieldalignment: struct with 56 pointer bytes could be 48 (govet)
```shell
service/ticketRobot/api/internal/config/config.go:59:17: fieldalignment: struct with 64 pointer bytes could be 56 (govet)
type DealerInfo struct {
                ^
}         
```

修复工具
```shell
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
```
```shell
fieldalignment -fix ./xxx/xxx
```
说明：
* -fix 会自动帮我们做到内存对齐，不需要我们做额外的操作；
* ./xxx/xxx 是我们的需要fix 的文件的包路径即pkg的目录

处理后的结构体
```go
// 处理人信息
type DealerInfo struct {
    DealId2DealerName map[string]string
    DealId2Weight     map[string]int64
+	ManagerId         string
+	FakePeople        string
+	DealerIdList      []string
}

```

## 参考资料
1. [Go语言详解内存对齐](https://blog.csdn.net/qq_53267860/article/details/124881698)