package main

import "fmt"

// 切片使用copy和等号复制的区别
//1.性能方面：copy复制会比等号复制慢。
//2.复制方式：copy复制为值复制，改变原切片的值不会影响新切片。而等号复制为指针复制，改变原切片或新切片都会对另一个产生影响。

// 如果你需要拷贝的对象中没有引用类型，那么对于Golang而言使用浅拷贝就可以了。
func main() {
	a := []int{0, 1, 2, 3, 4, 5, 6}
	fmt.Println(lastNumsBySlice(a))
	fmt.Println(lastNumsByCopy(a))

}

// 原始切片上操作，底层数组没有发生变化，内存一直占用，直到没有变量引用该数组，这种操作不推荐
func lastNumsBySlice(origin []int) (sliceInfo []int) {
	sliceInfo = origin[len(origin)-2:]
	return
}

// 推荐做法：copy
func lastNumsByCopy(origin []int) []int {
	result := make([]int, 2)
	num := copy(result, origin[len(origin)-2:])
	fmt.Printf("拷贝的数量有%v\n", num)
	return result

}
