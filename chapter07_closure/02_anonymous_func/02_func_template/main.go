package main

//  gin框架路由设计模式相同
/*	总结：
闭包的使用有点类似于面向对象设计模式中的模版模式，在模版模式中是在父类中定义公共的行为执行序列，
	然后子类通过重载父类的方法来实现特定的操作，而在Go语言中我们使用闭包实现了同样的效果
*/
import (
	"fmt"
)

type FilterFunc func(ele interface{}) interface{}

/*
  公共操作:对数据进行特殊操作
*/
func Data(arr interface{}, filterFunc FilterFunc) interface{} {

	slice := make([]int, 0)
	array, _ := arr.([]int)

	for _, value := range array {
		integer, ok := filterFunc(value).(int)
		if ok {
			slice = append(slice, integer)
		}

	}
	return slice
}

/*
  具体操作:奇数变偶数（这里可以不使用接口类型,直接使用int类型)
*/
func EvenFilter(ele interface{}) interface{} {

	integer, ok := ele.(int)
	if ok {
		if integer%2 == 1 {
			integer = integer + 1
		}
	}
	return integer
}

/*
  具体操作:偶数变奇数（这里可以不使用接口类型,直接使用int类型)
*/
func OddFilter(ele interface{}) interface{} {

	integer, ok := ele.(int)

	if ok {
		if integer%2 != 1 {
			integer = integer + 1
		}
	}

	return integer
}

func main() {
	sliceEven := make([]int, 0)
	sliceEven = append(sliceEven, 1, 2, 3, 4, 5)
	sliceEven = Data(sliceEven, EvenFilter).([]int)
	fmt.Println(sliceEven) //[2 2 4 4 6]

	sliceOdd := make([]int, 0)
	sliceOdd = append(sliceOdd, 1, 2, 3, 4, 5)
	sliceOdd = Data(sliceOdd, OddFilter).([]int)
	fmt.Println(sliceOdd) //[1 3 3 5 5]

}
