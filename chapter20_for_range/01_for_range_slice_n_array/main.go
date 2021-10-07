package main

import "fmt"

func main() {
	//Array()
	Slice()
	StructArray()

}

func Array() {
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int
	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r[i] = v
	}
	fmt.Println("r = ", r) // r =  [1 2 3 4 5]

	fmt.Println("a = ", a) // a =  [1 12 13 4 5]

}

func Slice() {
	var a = []int{1, 2, 3, 4, 5}
	var r = make([]int, len(a))
	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r[i] = v
	}
	fmt.Println("r = ", r) // r =  [1 12 13 4 5]

	fmt.Println("a = ", a) // a =  [1 12 13 4 5]

}

type T struct {
	n int
}

func StructArray() {

	ts := [2]T{}
	for i, t := range ts {
		switch i {
		case 0:
			t.n = 3
			ts[1].n = 9
		case 1:
			fmt.Println(t.n, " ") //0
		}
	}
	fmt.Println(ts) //[{0} {9}]
	// for-range 循环数组时使用的是数组 ts 的副本，所以 t.n = 3 的赋值操作不会影响原数组。
	//但 ts[1].n = 9这种方式操作的确是原数组的元素值，所以是会发生变化的。这也是我们推崇的方法。
}

/*
数组情况下
r 分析：
	对于所有的 range 循环 Go 语言都会在编译期为遍历对象创造一个副本，所以循环中通过短声明的变量修改值不会影响原循环数组的值。
	第一次遍历时修改了 a 的第二个和第三个元素，理论上第二次和第三次遍历时 r 应该能取到 a 修改后的值，
	但是我们刚说了 range 遍历开始前会创建副本，也就是说 range 的是 a 的副本而不是 a 本身。所以 r 赋值时用的都是 a 的副本的 value 值，所以不变。

a分析
	if 语句中赋值语句是用的 a[1],a[2] 这时候是真的修改 a 的值的，所以 a 变了，这里也是我们推荐的用法
*/
/*
切片情况下
分析：
	循环过程中依然创建了原切片的副本，但是因为切片自身的结构，创建的副本依然和原切片共享底层数组，只要没发生扩容，他们的值发生变化时就是同步变化的。效果就如同数组时range &a 一样了。


*/
