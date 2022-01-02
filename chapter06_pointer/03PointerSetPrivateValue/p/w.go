package p

import (
	"fmt"
	"unsafe"
)

type W struct {
	b byte
	i int32
	j int64
}

// b是byte类型，占1个字节；i是int32类型，占4个字节；j是int64类型，占8个字节: 1+4+8=13

func init() {
	fmt.Println("---初始化Init开始---")
	var w *W = new(W)
	fmt.Printf("Struct_W_size=%d\n", unsafe.Sizeof(*w))   // size=16
	fmt.Printf("w.b_alignSize=%d\n", unsafe.Alignof(w.b)) // alignSize=1
	fmt.Printf("w.i_alignSize=%d\n", unsafe.Sizeof(w.i))  // alignSize=4
	fmt.Printf("w.j_alignSize=%d\n", unsafe.Sizeof(w.j))  // alignSize=8
	fmt.Println("---初始化Init结束---")
}

// 13 !=16
// 原因：发生了对齐。
//		在struct中，它的对齐值是它的成员中的最大对齐值。每个成员类型都有它的对齐值，可以用unsafe.Alignof方法来计算，4
//		比如unsafe.Alignof(w.b)就可以得到b在w中的对齐值。
//		同理，我们可以计算出w.b的对齐值是1，w.i的对齐值是4，w.j的对齐值也是4

/* 分析：
对齐值最小是1，这是因为存储单元是以字节为单位。
所以b就在w的首地址，而i的对齐值是4，它的存储地址必须是4的倍数，因此，在b和i的中间有3个填充，同理j也需要对齐，
但因为i和j之间不需要填充，所以w的Sizeof值应该是13+3=16。

有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动GC。
当一个变量被移动，所有的保存改变量旧地址的指针必须同时被更新为变量移动后的新地址。
从垃圾收集器的视角来看，一个unsafe.Pointer是一个指向变量的指针，
因此当变量被移动是对应的指针也必须被更新；但是uintptr类型的临时变量只是一个普通的数字，所以其值不应该被改变。
*/
