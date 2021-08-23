package pkg

//var Id = 9527

var Name = "gopher"

// go tool compile -S pkg.go
/*
一。int类型
	go.cuinfo.packagename. SDWARFINFO dupok size=0
			0x0000 70 6b 67                                         pkg
	"".Id SNOPTRDATA size=8
			0x0000 37 25 00 00 00 00 00 00                          7%......
解释
	其中"".Id对应Id变量符号，变量的内存大小为8个字节。变量的初始化内容为37 25 00 00 00 00 00 00，对应十六进制格式的0x2537，
	对应十进制为9527。SNOPTRDATA是相关的标志，其中NOPTR表示数据中不包含指针数据。
	以上的内容只是目标文件对应的汇编，和Go汇编语言虽然相似当并不完全等价。

二。string类型
	go.cuinfo.packagename. SDWARFINFO dupok size=0
			0x0000 70 6b 67                                         pkg
	go.string."gopher" SRODATA dupok size=6
			0x0000 67 6f 70 68 65 72                                gopher
	"".Name SDATA size=16
			0x0000 00 00 00 00 00 00 00 00 06 00 00 00 00 00 00 00  ................
			rel 0+8 t=1 go.string."gopher"+0

解释
	因为Go语言的字符串并不是值类型，Go字符串其实是一种只读的引用类型。如果多个代码中出现了相同的"gopher"只读字符串时，程序链接后可以引用的同一个符号go.string."gopher"。
	因此，该符号有一个SRODATA标志表示这个数据在只读内存段，dupok表示出现多个相同标识符的数据时只保留一个就可以了
	真正的Go字符串变量Name对应的大小却只有16个字节了。其实Name变量并没有直接对应“gopher”字符串，而是对应16字节大小的reflect.StringHeader结构体
	type reflect.StringHeader struct {
		Data uintptr
		Len  int
	}前8个字节对应底层真实字符串数据的指针，也就是符号go.string."gopher"对应的地址。后8个字节对应底层真实字符串数据的有效长度，这里是6个字节

*/
