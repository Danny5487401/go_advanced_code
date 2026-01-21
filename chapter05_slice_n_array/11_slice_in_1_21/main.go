package main

import (
	"fmt"
	"slices"
)

func main() {
	// Contains 函数用于判断切片里是否包含指定元素
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(slices.Contains(nums, 3)) // true

	// Equal 函数用于比较两个切片是否相等，要求切片的元素类型必须是可比较(comparable)的
	numsCopy := []int{1, 2, 3, 4, 5}
	fmt.Println(slices.Equal(nums, numsCopy)) // true

	// Clip 函数用于删除切片中未使用的容量，执行操作后，切片的长度 = 切片的容量
	s := make([]int, 0, 8)
	s = append(s, 1, 2, 3, 4)
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // len: 4, cap: 8
	s = slices.Clip(s)
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // len: 4, cap: 4

	// Clone 函数返回一个拷贝的切片副本，元素是赋值复制，因此是浅拷贝. 由于是浅拷贝，修改副本切片里的元素，原切片的元素也会更新
	type User struct {
		Name string
	}
	users := []*User{
		&User{Name: "Danny"},
	}
	copiedSlice := slices.Clone(users)
	copiedSlice[0].Name = "Joy"
	fmt.Println(users[0].Name == copiedSlice[0].Name) // true

	// Compact 函数会将切片里连续的相同元素替换为一个元素。
	compactSlice := []int{1, 2, 2, 3, 3, 4, 5}
	newSlice := slices.Compact(compactSlice)
	fmt.Println(newSlice) // [1 2 3 4 5]

	// Delete 函数的功能是从指定切片 s 中删除指定范围 s[i:j] 的元素，并返回新的的切片。
	numbers := []int{1, 2, 3, 4, 5}
	newNumbers := slices.Delete(numbers, 1, 3)
	fmt.Println(newNumbers) // [1 4 5]

}
