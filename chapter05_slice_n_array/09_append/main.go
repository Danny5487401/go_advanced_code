package main

import "fmt"

func main() {
	//append1()
	//append2()
	append3()

}

/*
a[low : high : max]
//等价于=>
a[low: high]

0 <= low <= high <= max <= cap(a)
*/

func append1() {

	array := [...]int{1, 2, 3, 4, 5}
	s1 := array[:2]
	s2 := array[2:]
	s1 = append(s1, 999)
	fmt.Printf("array:%v len(array):%v cap(array):%v\n", array, len(array), cap(array)) // array:[1 2 999 4 5] len(array):5 cap(array):5
	fmt.Printf("s1:%v len(s1):%v cap(s1):%v\n", s1, len(s1), cap(s1))                   // s1:[1 2 999] len(s1):3 cap(s1):5
	fmt.Printf("s2:%v len(s2):%v cap(s2):%v\n", s2, len(s2), cap(s2))                   // s2:[999 4 5] len(s2):3 cap(s2):3

}

func append2() {
	array := [...]int{1, 2, 3, 4, 5}
	s1 := array[:2]
	s2 := array[2:]
	s1 = append(s1, 999, 888, 777)
	fmt.Printf("array:%v len(array):%v cap(array):%v\n", array, len(array), cap(array)) // array:[1 2 999 888 777] len(array):5 cap(array):5
	fmt.Printf("s1:%v len(s1):%v cap(s1):%v\n", s1, len(s1), cap(s1))                   //s1:[1 2 999 888 777] len(s1):5 cap(s1):5
	fmt.Printf("s2:%v len(s2):%v cap(s2):%v\n", s2, len(s2), cap(s2))                   //s2:[999 888 777] len(s2):3 cap(s2):3
}

// 如果切片append后的长度超越了容量，这时底层数组将会更换，且容量也会相应增加，原有的底层数组保持不变，原有的其他切片表达式也保持不变。
func append3() {
	array := [...]int{1, 2, 3, 4, 5}
	s1 := array[:2]
	s2 := array[2:]
	s1 = append(s1, 999, 888, 777, 666)
	fmt.Printf("array:%v len(array):%v cap(array):%v\n", array, len(array), cap(array)) // array:[1 2 3 4 5] len(array):5 cap(array):5
	fmt.Printf("s1:%v len(s1):%v cap(s1):%v\n", s1, len(s1), cap(s1))                   //s1:[1 2 999 888 777 666] len(s1):6 cap(s1):10
	fmt.Printf("s1:%v len(s1):%v cap(s1):%v\n", s1, len(s1), cap(s1))                   //s1:[1 2 999 888 777 666] len(s1):6 cap(s1):10
	fmt.Printf("s2:%v len(s2):%v cap(s2):%v\n", s2, len(s2), cap(s2))                   //
}
