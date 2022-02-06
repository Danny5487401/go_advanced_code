package main

import "fmt"

func main() {
	//MapAddElem()

	MapDelElem()
}

func MapAddElem() {
	var m = map[string]int{"A": 21,
		"B": 22,
		"C": 23,
	}
	counter := 0
	for k, v := range m {
		counter++
		fmt.Println(k, v)
		key := fmt.Sprintf("%s%d", "D", counter)
		m[key] = 24 //给map增加了新元素
	}
	fmt.Println("counter is ", counter) // counter is  7
	fmt.Println(m)                      // map[A:21 B:22 C:23 D1:24 D2:24 D3:24 D4:24 D5:24 D6:24 D7:24]
	/*
		A 21
		B 22
		C 23
		D1 24
		D2 24
		D3 24
		D4 24
		counter is  7
		map[A:21 B:22 C:23 D1:24 D2:24 D3:24 D4:24 D5:24 D6:24 D7:24]

	*/
}

func MapDelElem() {
	var m = map[string]int{"A": 21,
		"B": 22,
		"C": 23,
	}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "A")
		}
		counter++
		fmt.Println(k, v)

	}
	fmt.Println("counter is ", counter) //  2或者3
	/*
		for range map 是无序的，如果第一次循环到 A，则输出 3，否则输出 2。如果 map 中的元素在还没有被遍历到时就被移除了，后续的迭代中这个元素就不会再出现
	*/
}
