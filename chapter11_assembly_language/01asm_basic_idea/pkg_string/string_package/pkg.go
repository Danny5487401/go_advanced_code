package string_package

var Name string

/*
	其实Go汇编语言中定义的数据并没有所谓的类型，每个符号只不过是对应一块内存而已，因此NameData符号也是没有类型的。
	当Go语言的垃圾回收器在扫描到NameData变量的时候，无法知晓该变量内部是否包含指针，通过给NameData变量增加一个NOPTR标志，表示其中不会包含指针数据可以修复该错误：
*/
