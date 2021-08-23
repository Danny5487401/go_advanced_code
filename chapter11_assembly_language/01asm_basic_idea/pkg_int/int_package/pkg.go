package int_package

var Id int

/*
DATA symbol+offset(SB)/width, value
	其中symbol为变量在汇编语言中对应的标识符，offset是符号开始地址的偏移量，width是要初始化内存的宽度大小，value是要初始化的值。
	其中当前包中Go语言定义的符号symbol，在汇编代码中对应·symbol，其中“·”中点符号为一个特殊的unicode符号


pkg_amd64.s的后缀名表示AMD64环境下的汇编代码文件
*/
